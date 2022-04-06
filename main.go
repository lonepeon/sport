package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3" // sqlite3 adapter

	"github.com/lonepeon/golib/env"
	"github.com/lonepeon/golib/job"
	"github.com/lonepeon/golib/logger"
	"github.com/lonepeon/golib/sqlutil"
	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/golib/web/authenticationstore"
	"github.com/lonepeon/golib/web/sessionstore"
	"github.com/lonepeon/sport/internal/application/service"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/infrastructure/annotation"
	"github.com/lonepeon/sport/internal/infrastructure/gpx"
	domainjob "github.com/lonepeon/sport/internal/infrastructure/job"
	"github.com/lonepeon/sport/internal/infrastructure/mapbox"
	"github.com/lonepeon/sport/internal/infrastructure/s3"
	"github.com/lonepeon/sport/internal/infrastructure/sqlite"
	"github.com/lonepeon/sport/internal/infrastructure/www"
	"github.com/lonepeon/sport/internal/repository"
)

type Repository struct {
	*s3.Bucket
	sqlite.SQLite
	*mapbox.Mapbox
	gpx.GPX
	annotation.Annotation
}

type Config struct {
	SQLitePath         string   `env:"SPORT_SQLITE_PATH,default=./sport.sqlite"`
	SessionKey         string   `env:"SPORT_SESSION_KEY,required=true"`
	UploadFolder       string   `env:"SPORT_UPLOAD_FOLDER,default=./tmp/uploads,required=true"`
	WebAddress         string   `env:"SPORT_WEB_ADDR,required=true"`
	CDNURL             string   `env:"SPORT_CDN_URL,required=true"`
	AWSAccessKeyID     string   `env:"SPORT_AWS_ACCESS_KEY_ID,required=true"`
	AWSSecretAccessKey string   `env:"SPORT_AWS_SECRET_ACCESS_KEY,required=true"`
	AWSRegion          string   `env:"SPORT_AWS_REGION,required=true"`
	AWSBucket          string   `env:"SPORT_AWS_BUCKET,required=true"`
	AWSEndpointURL     string   `env:"SPORT_AWS_ENDPOINT_URL"`
	MapboxEndpointURL  string   `env:"SPORT_MAPBOX_ENDPOINT_URL"`
	MapboxToken        string   `env:"SPORT_MAPBOX_TOKEN,required=true"`
	Users              []string `env:"SPORT_USERS,required=true,sep=;"`
}

//go:embed templates/*
var htmlTemplateFS embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	log, closer := logger.NewLogger(os.Stdout)
	defer func() {
		if err := closer(); err != nil {
			fmt.Fprintf(os.Stderr, "can't flush logs: %v", err)
		}
	}()

	log = log.WithFields(logger.String("app-version", "HEAD"), logger.String("app-name", "sport"))

	var cfg Config
	if err := env.Load(&cfg); err != nil {
		return fmt.Errorf("can't load config: %v", err)
	}

	db, err := initDatabase(log, cfg.SQLitePath)
	if err != nil {
		return fmt.Errorf("can't initialize database: %v", err)
	}

	sessionstore := sessionstore.NewSQLite(db, sessions.Options{
		HttpOnly: true,
		MaxAge:   1 * 60 * 60 * 24 * 2,
	}, []byte(cfg.SessionKey))

	repo := repository.NewLogger(log, Repository{
		Bucket: initBucket(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			cfg.AWSRegion,
			cfg.AWSBucket,
			cfg.AWSEndpointURL,
		),
		SQLite: sqlite.New(db),
		Mapbox: initMapbox(cfg.MapboxToken, cfg.MapboxEndpointURL),
	})

	application := service.NewApplication(repo)

	jobServer, jobClient := initJob(db, log,
		domainjob.NewTrackRunningSessionJob(application),
		domainjob.NewDeleteRunningSessionJob(application),
	)

	auth, err := initAutenticationMiddleware(sessionstore, cfg.Users)
	if err != nil {
		return err
	}
	webServer := initWebServer(log, sessionstore, cfg.CDNURL)
	webServer.HandleFunc("GET", "/login", auth.ShowLoginPage("/running-session/new"))
	webServer.HandleFunc("POST", "/login", auth.Login("/running-session/new"))
	webServer.HandleFunc("GET", "/logout", auth.Logout("/"))
	webServer.HandleFunc("GET", "/", auth.IdentifyCurrentUser(www.RunningSessionsIndex(application)))
	webServer.HandleFunc("GET", "/running-session/new", auth.EnsureAuthentication("/login", www.RunningSessionNew()))
	webServer.HandleFunc("POST", "/running-session", auth.EnsureAuthentication("/login", www.RunningSessionPost(jobClient, cfg.UploadFolder)))
	webServer.HandleFunc("GET", "/running-session/{slug}", auth.IdentifyCurrentUser((www.RunningSessionsShow(application))))
	webServer.HandleFunc("POST", "/running-session/{slug}/delete", auth.EnsureAuthentication("/login", www.RunningSessionsDelete(application, jobClient)))

	return waitForServersShutdown(log, jobServer, webServer, cfg.WebAddress)
}

func registerUsers(authBackend web.AuthenticationBackendStorer, rawUsers []string) error {
	for _, rawUserLine := range rawUsers {
		rawUser := strings.Split(rawUserLine, ":")
		if len(rawUser) != 2 {
			return fmt.Errorf("entry is expected to be base64(username):base64(password) but is '%s'", rawUserLine)
		}

		username, err := base64.URLEncoding.DecodeString(rawUser[0])
		if err != nil {
			return fmt.Errorf("username part is not a valid base64 value (value='%s')", rawUser[0])
		}
		password, err := base64.URLEncoding.DecodeString(rawUser[1])
		if err != nil {
			return fmt.Errorf("password part is not a valid base64 value (value='%s')", rawUser[1])
		}

		if _, err := authBackend.Register(string(username), string(password)); err != nil {
			return fmt.Errorf("can't register user (username='%s')", rawUser[0])
		}
	}

	return nil
}

func initDatabase(log *logger.Logger, sqlitePath string) (*sql.DB, error) {
	log.Infof("initialize database from %v", sqlitePath)
	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("can't open  sqlite file: %v", err)
	}

	jobMigrationsVersions, err := sqlutil.ExecuteMigrations(context.Background(), db, job.Migrations())
	if err != nil {
		return nil, fmt.Errorf("can't run job migrations: %v", err)
	}
	log.Infof("database executed new sql job migrations %s", strings.Join(jobMigrationsVersions, ", "))

	trackerMigrationsVersions, err := sqlutil.ExecuteMigrations(context.Background(), db, sqlite.Migrations())
	if err != nil {
		return nil, fmt.Errorf("can't run tracker migrations: %v", err)
	}
	log.Infof("database executed new sql application migrations %s", strings.Join(trackerMigrationsVersions, ", "))

	sessionStoreMigrationsVersions, err := sqlutil.ExecuteMigrations(context.Background(), db, sessionstore.Migrations())
	if err != nil {
		return nil, fmt.Errorf("can't run session store migrations: %v", err)
	}
	log.Infof("database executed new sql session store migrations %s", strings.Join(sessionStoreMigrationsVersions, ", "))

	return db, nil
}

func initJob(db *sql.DB, log *logger.Logger, jobHandlers ...job.Handler) (*job.Server, *job.Client) {
	reg := job.NewRegistry()
	for _, jobHandler := range jobHandlers {
		reg.Register(jobHandler)
	}
	jobServer := job.NewServer(db, reg, log)
	jobClient := jobServer.Client()

	return jobServer, jobClient
}

func initBucket(accessKeyID string, secretAccessKey string, region string, bucketName string, endpointURL string) *s3.Bucket {
	os.Setenv("AWS_ACCESS_KEY_ID", accessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretAccessKey)

	s3Bucket := s3.NewBucket(bucketName, region)
	if endpointURL != "" {
		s3Bucket.Endpoint = endpointURL
	}

	return s3Bucket
}

func initWebServer(log *logger.Logger, sessionstore sessions.Store, cdnURL string) *web.Server {
	tmpl := web.TmplConfiguration{
		FS:                          htmlTemplateFS,
		Layout:                      "templates/layout.html.tmpl",
		ErrorLayout:                 "templates/layout.html.tmpl",
		RedirectionTemplate:         "templates/30x.html.tmpl",
		NotFoundTemplate:            "templates/404.html.tmpl",
		InternalServerErrorTemplate: "templates/500.html.tmpl",
		UnauthorizedTemplate:        "templates/401.html.tmpl",
	}

	// TODO: add form with CSRF
	webServer := web.NewServer(log, tmpl, sessionstore)
	webServer.AddTemplateFuncs(template.FuncMap{
		"fmtdatetime": func(dt time.Time) string {
			return dt.Format("2006/01/02 15:04")
		},
		"modulo": func(value int, base int) bool {
			return value%base == 0
		},
		"ternary": func(iftrue interface{}, iffalse interface{}, cond bool) interface{} {
			if cond {
				return iftrue
			}

			return iffalse
		},
		"mapurl": func(fname domain.MapFilePath) string {
			return cdnURL + "/" + string(fname)
		},
		"shareablemapurl": func(fname domain.ShareableMapFilePath) string {
			return cdnURL + "/" + string(fname)
		},
	})

	return webServer
}

func spawnJobServer(log *logger.Logger, server *job.Server, reporter chan<- error) (func(), func(time.Duration) error) {
	start := func() {
		log.Info("starting job server")
		if err := server.ListenAndServe(); err != nil {
			reporter <- fmt.Errorf("job server failed: %v", err)
		}
	}

	shutdown := func(d time.Duration) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	}

	return start, shutdown
}

func spawnWebServer(log *logger.Logger, server *web.Server, webAddress string, reporter chan<- error) (func(), func(time.Duration) error) {
	start := func() {
		log.Infof("starting web server on %s", webAddress)
		if err := server.ListenAndServe(webAddress); err != nil {
			reporter <- fmt.Errorf("web server failed: %v", err)
		}
	}

	stop := func(d time.Duration) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	}

	return start, stop
}

func waitForServersShutdown(log *logger.Logger, jobServer *job.Server, webServer *web.Server, webAddress string) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	webErr := make(chan error)
	jobErr := make(chan error)

	jobServerStarter, jobServerStopper := spawnJobServer(log, jobServer, jobErr)
	webServerStarter, webServerStopper := spawnWebServer(log, webServer, webAddress, webErr)

	go jobServerStarter()
	go webServerStarter()

	var err error
	shutdownDuration := 20 * time.Second

	select {
	case <-sigs:
		err = jobServerStopper(shutdownDuration)
		if webErr := webServerStopper(shutdownDuration); err != nil {
			err = fmt.Errorf("%v; %v", err, webErr)
		}
	case e := <-webErr:
		log.Error(e.Error())
		err = jobServerStopper(shutdownDuration)
	case e := <-jobErr:
		log.Error(e.Error())
		err = webServerStopper(shutdownDuration)
	}

	return err
}

func initMapbox(token string, endpointURL string) *mapbox.Mapbox {
	box := mapbox.New(token)
	if endpointURL != "" {
		box.EndpointURL = endpointURL
	}

	return box
}

func initAutenticationMiddleware(store sessions.Store, users []string) (web.Authentication, error) {
	authenticationBrowserStore := web.NewCurrentAuthenticatedUserSessionStore(store)
	authenticationBackendstore := authenticationstore.NewInMemory()
	if err := registerUsers(authenticationBackendstore, users); err != nil {
		return web.Authentication{}, fmt.Errorf("can't parse USERS environment variable: %v", err)
	}

	auth := web.NewAuthentication(authenticationBrowserStore, authenticationBackendstore, "templates/login/new.html.tmpl")

	return auth, nil
}
