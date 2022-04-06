package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var ErrUserInvalidCredentials = errors.New("invalid user credentials")
var ErrUserNotFound = errors.New("current user not found")
var ErrUserAlreadyExist = errors.New("user already registered")

type AuthenticationUserID string

func (a AuthenticationUserID) String() string {
	return string(a)
}

type AuthenticationUser struct {
	ID       AuthenticationUserID
	Username string
}

type CurrentAuthenticatedUserSessionStore struct {
	store sessions.Store
}

func NewCurrentAuthenticatedUserSessionStore(store sessions.Store) CurrentAuthenticatedUserSessionStore {
	return CurrentAuthenticatedUserSessionStore{
		store: store,
	}
}

func (c CurrentAuthenticatedUserSessionStore) Clear(w http.ResponseWriter, r *http.Request) error {
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

func (c CurrentAuthenticatedUserSessionStore) StoreUserID(w http.ResponseWriter, r *http.Request, id AuthenticationUserID) error {
	session, err := c.store.Get(r, "auth")
	if err != nil {
		return fmt.Errorf("can't get session 'auth': %v", err)
	}

	session.Values["auth_id"] = id.String()

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("can't save username to session 'auth': %v", err)
	}

	return nil
}

func (c CurrentAuthenticatedUserSessionStore) CurrentUserID(r *http.Request) (AuthenticationUserID, error) {
	session, err := c.store.Get(r, "auth")
	if err != nil {
		return "", fmt.Errorf("can't get session 'auth': %v", err)
	}

	id, ok := session.Values["auth_id"].(string)
	if !ok {
		return "", nil
	}

	return AuthenticationUserID(id), nil
}

type AuthenticationFrontendStorer interface {
	StoreUserID(http.ResponseWriter, *http.Request, AuthenticationUserID) error
	Clear(http.ResponseWriter, *http.Request) error
	// CurrentUserID returns:
	// - the current authenticated user username if authenticated
	// - an empty string when there are no authenticated user
	// - an error if something went wrong during the retrieval process
	//
	// We don't expect an error if no users are currently logged in, just an empty string
	CurrentUserID(*http.Request) (AuthenticationUserID, error)
}

type AuthenticationBackendStorer interface {
	// Register returns the ID of the register user
	// It returns an ErrUserAlreadyExist if the username is already taken
	Register(username string, password string) (AuthenticationUserID, error)

	// Authenticate returns the ID of the user if the authentication is successful
	// It returns an ErrUserInvalidCredentials if the username or password is invalid
	Authenticate(username string, password string) (AuthenticationUserID, error)

	// Lookup returns the user from its ID
	// It returns an ErrUserNotFound if the ID doesn't belong to any user
	Lookup(id AuthenticationUserID) (AuthenticationUser, error)
}

type Authentication struct {
	frontendStorage   AuthenticationFrontendStorer
	backendStorage    AuthenticationBackendStorer
	loginTemplatePath string
}

func NewAuthentication(front AuthenticationFrontendStorer, back AuthenticationBackendStorer, loginTemplatePath string) Authentication {
	return Authentication{
		frontendStorage:   front,
		backendStorage:    back,
		loginTemplatePath: loginTemplatePath,
	}
}

func (a Authentication) ShowLoginPage(alreadyLoggedinRedirectPath string) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		userID, err := a.frontendStorage.CurrentUserID(r)
		if err != nil {
			return ctx.Response(200, a.loginTemplatePath, nil)
		}

		if userID != "" {
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

		user, err := a.authenticateUser(username, password)
		if err != nil {
			if errors.Is(err, ErrUserInvalidCredentials) {
				ctx.AddFlash(NewFlashMessageError("username/password combination is invalid"))
				return ctx.Response(http.StatusOK, a.loginTemplatePath, nil)
			}
			ctx.AddFlash(NewFlashMessageError("something wrong happened. Please try again."))
			response := ctx.Response(http.StatusOK, a.loginTemplatePath, nil)
			response.LogMessage = fmt.Sprintf("can't authenticate user because of storage error: %v", err)
			return response
		}

		if err := a.frontendStorage.StoreUserID(w, r, user.ID); err != nil {
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
		if err := a.frontendStorage.Clear(w, r); err != nil {
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
		user, err := a.loadCurrentUser(r)
		if err != nil {
			user = AuthenticationUser{}
		}

		a.setCtxData(ctx, user)

		return h(ctx, w, r)
	}
}

func (a Authentication) EnsureAuthentication(loginRedirectPath string, h HandlerFunc) HandlerFunc {
	return func(ctx Context, w http.ResponseWriter, r *http.Request) Response {
		user, err := a.loadCurrentUser(r)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				redirection := ctx.Redirect(w, http.StatusFound, loginRedirectPath)
				redirection.LogMessage = fmt.Sprintf("redirecting to %v: current user is not authenticated", loginRedirectPath)
				return redirection

			}

			redirection := ctx.Redirect(w, http.StatusFound, loginRedirectPath)
			redirection.LogMessage = fmt.Sprintf("redirecting to %v: can't load current user: %v", loginRedirectPath, err)
			return redirection
		}

		a.setCtxData(ctx, user)
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

func (a Authentication) setCtxData(ctx Context, user AuthenticationUser) {
	ctx.AddData("Authentication", map[string]interface{}{
		"IsLoggedIn": user.ID != "",
		"Username":   user.Username,
		"User":       user,
	})
}

func (a Authentication) authenticateUser(username string, password string) (AuthenticationUser, error) {
	id, err := a.backendStorage.Authenticate(username, password)
	if err != nil {
		if errors.Is(err, ErrUserInvalidCredentials) {
			return AuthenticationUser{}, fmt.Errorf("cannot authenticate user (username=%s): %w", username, err)
		}

		return AuthenticationUser{}, fmt.Errorf("unexpected error while authenticating user (username=%s): %v", username, err)
	}

	return a.backendStorage.Lookup(id)
}

func (a Authentication) loadCurrentUser(r *http.Request) (AuthenticationUser, error) {
	id, err := a.frontendStorage.CurrentUserID(r)
	if err != nil {
		return AuthenticationUser{}, fmt.Errorf("can't retrieve user in session: %v", err)
	}

	if id == "" {
		return AuthenticationUser{}, fmt.Errorf("no user in session: %w", ErrUserNotFound)
	}

	return a.backendStorage.Lookup(id)
}
