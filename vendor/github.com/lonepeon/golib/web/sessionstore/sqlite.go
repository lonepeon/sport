package sessionstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const (
	createdAtField = "createdAt"
	expiredAtField = "expiredAt"
)

//go:generate go run ../../sqlutil/cmd/sql-migration ./sqlite_scripts

type Session struct {
	ID        string
	Data      string
	CreatedAt string
	ExpiredAt string
}

type SQLite struct {
	db *sql.DB

	Codecs  []securecookie.Codec
	Options sessions.Options
}

func NewSQLite(db *sql.DB, options sessions.Options, keyPairs ...[]byte) *SQLite {
	return &SQLite{db: db, Options: options, Codecs: securecookie.CodecsFromPairs(keyPairs...)}
}

// Get implements gorilla.sessions: an error if using the Registry infrastructure to cache the session.
func (s *SQLite) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New implements gorilla.sessions: should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (s *SQLite) New(r *http.Request, name string) (*sessions.Session, error) {
	sess := sessions.NewSession(s, name)

	option := s.Options
	sess.Options = &option
	sess.IsNew = true

	cookie, err := r.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return sess, nil
		}
		return sess, fmt.Errorf("can't extract cookie from request (cookie=%s): %w", name, err)
	}

	if err = securecookie.DecodeMulti(name, cookie.Value, &sess.ID, s.Codecs...); err != nil {
		return sess, fmt.Errorf("can't decode cookie value: %w", err)
	}

	if err := s.fillSession(r.Context(), sess); err != nil {
		return sess, fmt.Errorf("cant fill session: %w", err)
	}

	sess.IsNew = false

	return sess, nil
}

// Save implements gorilla.sessions: should persist session to the underlying store implementation.
func (s *SQLite) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	if err := s.saveSession(r.Context(), sess); err != nil {
		return err
	}

	data, err := securecookie.EncodeMulti(sess.Name(), sess.ID, s.Codecs...)
	if err != nil {
		return fmt.Errorf("can't encode session ID (id=%s): %w", sess.ID, err)
	}

	http.SetCookie(w, sessions.NewCookie(sess.Name(), data, sess.Options))

	return nil
}

func (s *SQLite) saveSession(ctx context.Context, sess *sessions.Session) error {
	if sess.ID == "" || sess.IsNew {
		return s.insertSession(ctx, sess)
	}

	return s.updateSession(ctx, sess)
}

func (s *SQLite) insertSession(ctx context.Context, sess *sessions.Session) error {
	createdAt, ok := sess.Values[createdAtField]
	if !ok {
		createdAt = time.Now().Format(time.RFC3339)
	}

	expiredAt, ok := sess.Values[expiredAtField]
	if !ok {
		expiredAt = time.Now().Add(time.Second * time.Duration(sess.Options.MaxAge)).Format(time.RFC3339)
	}

	data, err := securecookie.EncodeMulti(sess.Name(), sess.Values, s.Codecs...)
	if err != nil {
		return fmt.Errorf("can't encode new session values: %w", err)
	}

	sessID := uuid.NewString()
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO sessions (id, data, created_at, expired_at) VALUES (?, ?, ?, ?)`,
		sessID, data, createdAt, expiredAt,
	)

	if err != nil {
		return fmt.Errorf("can't persist new session in database: %v", err)
	}

	sess.IsNew = false
	sess.ID = sessID

	return nil
}

func (s *SQLite) updateSession(ctx context.Context, sess *sessions.Session) error {
	data, err := securecookie.EncodeMulti(sess.Name(), sess.Values, s.Codecs...)
	if err != nil {
		return fmt.Errorf("can't encode new session values: %w", err)
	}

	result, err := s.db.ExecContext(ctx, `UPDATE sessions SET data = ? WHERE id = ?`, data, sess.ID)
	if err != nil {
		return fmt.Errorf("can't persist existing session in database (id=%s): %v", sess.ID, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't be sure if session has been updated: can't get the number of row affected: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("can't update session (id=%s): session not found", sess.ID)
	}

	return nil
}

func (s *SQLite) fillSession(ctx context.Context, sess *sessions.Session) error {
	row, err := s.db.QueryContext(ctx, "SELECT id, data, created_at, expired_at FROM sessions WHERE id = ?", sess.ID)
	if err != nil {
		return fmt.Errorf("can't query session table (id=%s): %v", sess.ID, err)
	}

	defer row.Close()

	if !row.Next() {
		return fmt.Errorf("can't find session with id (id=%s)", sess.ID)
	}

	var dto Session
	if err := row.Scan(&dto.ID, &dto.Data, &dto.CreatedAt, &dto.ExpiredAt); err != nil {
		return fmt.Errorf("can't scan session row (id=%s): %v", sess.ID, err)
	}

	expiredAt, err := time.Parse(time.RFC3339, dto.ExpiredAt)
	if err != nil {
		return fmt.Errorf("can't parse expired to a time (value=%s): %v", dto.ExpiredAt, err)
	}

	if expiredAt.Before(time.Now()) {
		return fmt.Errorf("can't get expired session (id=%s, expiredAt=%s)", sess.ID, dto.ExpiredAt)
	}

	if err := securecookie.DecodeMulti(sess.Name(), dto.Data, &sess.Values, s.Codecs...); err != nil {
		return fmt.Errorf("can't decode session data (id=%s): %w", sess.ID, err)
	}

	sess.Options.MaxAge = int(time.Since(expiredAt).Seconds())
	sess.Values[createdAtField] = dto.CreatedAt
	sess.Values[expiredAtField] = dto.ExpiredAt

	return nil
}

var _ sessions.Store = &SQLite{}
