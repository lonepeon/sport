package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type CurrentAuthenticatedUserFilesystemSessionStore struct {
	store sessions.Store
}

func NewCurrentAuthenticatedUserSessionStore(store sessions.Store) CurrentAuthenticatedUserFilesystemSessionStore {
	return CurrentAuthenticatedUserFilesystemSessionStore{
		store: store,
	}
}

func (c CurrentAuthenticatedUserFilesystemSessionStore) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := c.store.Get(r, "auth")
	if err != nil {
		return fmt.Errorf("can't get session 'auth': %v", err)
	}

	session.Values = nil

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("can't save cleared 'auth' session: %v", err)
	}

	return nil
}

func (c CurrentAuthenticatedUserFilesystemSessionStore) AuthenticateUsername(w http.ResponseWriter, r *http.Request, name string) error {
	session, err := c.store.Get(r, "auth")
	if err != nil {
		return fmt.Errorf("can't get session 'auth': %v", err)
	}

	session.Values["auth_username"] = name

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("can't save username to session 'auth': %v", err)
	}

	return nil
}

func (c CurrentAuthenticatedUserFilesystemSessionStore) CurrentUsername(r *http.Request) (string, error) {
	session, err := c.store.Get(r, "auth")
	if err != nil {
		return "", fmt.Errorf("can't get session 'auth': %v", err)
	}

	username, ok := session.Values["auth_username"].(string)
	if !ok {
		return "", nil
	}

	return username, nil
}

type CurrentAuthenticatedUserStorage interface {
	AuthenticateUsername(http.ResponseWriter, *http.Request, string) error
	Clear(http.ResponseWriter, *http.Request) error
	// CurrentUsername returns:
	// - the current authenticated user username if authenticated
	// - an empty string when there are no authenticated user
	// - an error if something went wrong during the retrieval process
	//
	// We don't expect an error if no users are currently logged in, just an empty string
	CurrentUsername(*http.Request) (string, error)
}

type AuthenticationUser struct {
	Username string
	Password string
}

type Authentication struct {
	storage           CurrentAuthenticatedUserStorage
	users             []AuthenticationUser
	loginTemplatePath string
}

func NewAuthentication(storage CurrentAuthenticatedUserStorage, users []AuthenticationUser, loginTemplatePath string) Authentication {
	return Authentication{
		storage:           storage,
		users:             users,
		loginTemplatePath: loginTemplatePath,
	}
}

func (a Authentication) ShowLoginPage(alreadyLoggedinRedirectPath string) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		username, err := a.storage.CurrentUsername(r)
		if err != nil {
			return ctx.Response(200, a.loginTemplatePath, nil)
		}

		if username != "" {
			dest := a.getRedirectionPath(r, alreadyLoggedinRedirectPath)
			return ctx.Redirect(w, http.StatusFound, dest)
		}

		return ctx.Response(200, a.loginTemplatePath, nil)
	}
}

func (a Authentication) Login(successfulLoginRedirectPath string) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		if err := r.ParseForm(); err != nil {
			ctx.AddFlash(NewFlashMessageError("can't parse request parameters. Please try again."))
			response := ctx.Response(http.StatusOK, a.loginTemplatePath, nil)
			response.LogMessage = fmt.Sprintf("can't parse form: %v", err)
			return response
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			ctx.AddFlash(NewFlashMessageError("username/password combination is required"))
			return ctx.Response(http.StatusOK, a.loginTemplatePath, nil)
		}

		user, isAuthorized := a.findUser(username, password)
		if !isAuthorized {
			ctx.AddFlash(NewFlashMessageError("username/password combination is invalid"))
			return ctx.Response(http.StatusOK, a.loginTemplatePath, nil)
		}

		if err := a.storage.AuthenticateUsername(w, r, user.Username); err != nil {
			ctx.AddFlash(NewFlashMessageError("something wrong happened. Please try again."))
			response := ctx.Response(http.StatusOK, a.loginTemplatePath, nil)
			response.LogMessage = fmt.Sprintf("can't login username because of storage error: %v", err)
			return response
		}

		dest := a.getRedirectionPath(r, successfulLoginRedirectPath)
		return ctx.Redirect(w, http.StatusFound, dest)
	}
}

func (a Authentication) Logout(successfulLogoutRedirectPath string) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		dest := a.getRedirectionPath(r, successfulLogoutRedirectPath)

		logMsg := fmt.Sprintf("redirecting to %v", dest)
		if err := a.storage.Clear(w, r); err != nil {
			ctx.AddFlash(NewFlashMessageError("We can't log you out. Please retry"))
			logMsg = fmt.Sprintf("%s: can't remove user from storage: %v", logMsg, err)
		}

		redirection := ctx.Redirect(w, http.StatusFound, dest)
		redirection.LogMessage = logMsg
		return redirection
	}
}

func (a Authentication) IdentifyCurrentUser(h HandlerFunc) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		username, err := a.storage.CurrentUsername(r)
		if err != nil {
			username = ""
		}

		a.setCtxData(ctx, username)

		return h(ctx, w, r)
	}
}

func (a Authentication) EnsureAuthentication(loginRedirectPath string, h HandlerFunc) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		username, err := a.storage.CurrentUsername(r)
		if err != nil {
			redirection := ctx.Redirect(w, http.StatusFound, loginRedirectPath)
			redirection.LogMessage = fmt.Sprintf("redirecting to %v: can't get current username: %v", loginRedirectPath, err)
			return redirection
		}

		if username == "" {
			redirection := ctx.Redirect(w, http.StatusFound, loginRedirectPath)
			redirection.LogMessage = fmt.Sprintf("redirecting to %v: current user is not authenticated", loginRedirectPath)
			return redirection
		}

		a.setCtxData(ctx, username)
		return h(ctx, w, r)
	}
}

func (a Authentication) getRedirectionPath(r *http.Request, defaultDest string) string {
	if err := r.ParseForm(); err != nil {
		return defaultDest
	}

	dest := r.Form.Get("to")
	if dest == "" {
		return defaultDest
	}

	return dest
}

func (a Authentication) setCtxData(ctx Context, username string) {
	ctx.AddData("Authentication", map[string]interface{}{
		"IsLoggedIn": username != "",
		"Username":   username,
	})
}

func (a Authentication) findUser(username string, password string) (AuthenticationUser, bool) {
	for _, authorizedUser := range a.users {
		if authorizedUser.Username == username && authorizedUser.Password == password {
			return authorizedUser, true
		}
	}

	return AuthenticationUser{}, false
}
