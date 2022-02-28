package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/lonepeon/golib/logger"
)

type TmplResponse struct {
	Flashes []interface{}
	Data    interface{}
}

type TmplConfiguration struct {
	FS                          embed.FS
	Layout                      string
	ErrorLayout                 string
	RedirectionTemplate         string
	NotFoundTemplate            string
	InternalServerErrorTemplate string
	UnauthorizedTemplate        string
}

type Server struct {
	logger       *logger.Logger
	tmplCfg      TmplConfiguration
	router       *mux.Router
	server       http.Server
	sessionStore *sessions.FilesystemStore
	tmplFuncs    template.FuncMap
}

func NewServer(log *logger.Logger, tmplCfg TmplConfiguration, sessionStore *sessions.FilesystemStore) *Server {
	return &Server{
		logger:  log,
		tmplCfg: tmplCfg,
		router:  mux.NewRouter(),
		server: http.Server{
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      45 * time.Second,
		},
		sessionStore: sessionStore,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	s.server.Addr = addr
	s.server.Handler = s.router
	return s.server.ListenAndServe()
}

func (s *Server) AddTemplateFuncs(funcs template.FuncMap) {
	s.tmplFuncs = funcs
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) HandleFunc(method string, urlpath string, h HandlerFunc) {
	handler := s.wrapRequest(method, urlpath, h)
	s.router.HandleFunc(urlpath, handler).Methods(method)
}

func (s *Server) Handle(method string, urlpath string, h Handler) {
	s.HandleFunc(urlpath, urlpath, h.Handle)
}

func (s *Server) wrapRequest(method string, urlpath string, h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.NewString()
		urlScheme := "http"
		if r.TLS != nil {
			urlScheme = "https"
		}
		reqLogger := s.logger.WithFields(
			logger.String("trace-id", traceID),
			logger.String("http.method", method),
			logger.String("http.url", r.URL.String()),
			logger.String("http.target", r.URL.RequestURI()),
			logger.String("http.host", r.Host),
			logger.String("http.scheme", urlScheme),
			logger.String("http.flavor", fmt.Sprintf("%d.%d", r.ProtoMajor, r.ProtoMinor)),
			logger.String("http.user_agent", r.UserAgent()),
		)

		session, err := s.sessionStore.Get(r, "trax")
		if err != nil {
			reqLogger.
				WithFields(logger.Int("http.status_code", http.StatusInternalServerError)).
				Errorf("can't get session store: %v", err)
			http.Error(w, "something wrong happened", http.StatusInternalServerError)
			return
		}

		ctx := ContextImpl{
			Context:           r.Context(),
			tmplConfiguration: s.tmplCfg,
			session:           session,
		}

		w.Header().Add("Trace-ID", traceID)

		httpCode, msg := s.writeResponse(&ctx, w, r, session, h(&ctx, w, r))

		reqLogger = reqLogger.WithFields(logger.Int("http.status_code", httpCode))
		logFn := reqLogger.Info
		if httpCode >= http.StatusInternalServerError {
			logFn = reqLogger.Error
		}

		logFn(msg)
	}
}

func (s *Server) writeResponse(ctx *ContextImpl, w http.ResponseWriter, r *http.Request, session *sessions.Session, resp Response) (int, string) {
	tmplName := path.Base(resp.Layout)
	if tmplName == "." {
		tmplName = path.Base(resp.Template)
	}

	tmpl, err := template.
		New(tmplName).
		Funcs(s.tmplFuncs).
		ParseFS(ctx.tmplConfiguration.FS, resp.Templates()...)

	if err != nil {
		return s.write500(w, s.wrapLogMessage(resp.LogMessage, "can't parse templates from %s: %v", strings.Join(resp.Templates(), ","), err))

	}

	tmplResponse := TmplResponse{
		Data: resp.Data,
	}

	if resp.HTTPCode < 300 || resp.HTTPCode >= 400 {
		tmplResponse.Flashes = session.Flashes()
	}

	if err := session.Save(r, w); err != nil {
		return s.write500(w, s.wrapLogMessage(resp.LogMessage, "can't save session: %v", err))
	}

	w.WriteHeader(resp.HTTPCode)
	if err := tmpl.Execute(w, tmplResponse); err != nil {
		return s.write500(w, s.wrapLogMessage(resp.LogMessage, "can't execute templates from %s: %v", strings.Join(resp.Templates(), ","), err))
	}

	return resp.HTTPCode, resp.LogMessage
}

func (s *Server) write500(w http.ResponseWriter, err error) (int, string) {
	http.Error(w, "something wrong happened", http.StatusInternalServerError)
	return http.StatusInternalServerError, err.Error()
}

func (s *Server) wrapLogMessage(original string, format string, args ...interface{}) error {
	return fmt.Errorf("%v. original log: %v", fmt.Sprintf(format, args...), original)
}
