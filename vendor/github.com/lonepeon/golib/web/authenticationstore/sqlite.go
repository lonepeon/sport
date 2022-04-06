package authenticationstore

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/lonepeon/golib/sqlutil/sqliteutil"
	"github.com/lonepeon/golib/web"
)

//go:generate go run ../../sqlutil/cmd/sql-migration ./sqlite_scripts

type SQLite struct {
	db     *sql.DB
	pepper string

	HashCost int
}

func NewSQLite(db *sql.DB, pepper string) *SQLite {
	return &SQLite{db: db, pepper: pepper, HashCost: bcrypt.DefaultCost}
}

func (s *SQLite) Authenticate(username string, password string) (web.AuthenticationUserID, error) {
	rows, err := s.db.Query("SELECT id, password, salt FROM authentication_user WHERE username = ?", username)
	if err != nil {
		return "", fmt.Errorf("can't execute lookup query: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("%w: invalid username", web.ErrUserInvalidCredentials)
	}

	var id, hashedPassword, salt string
	if err := rows.Scan(&id, &hashedPassword, &salt); err != nil {
		return "", fmt.Errorf("can't scan SQL columns: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt+s.pepper))
	if err != nil {
		return "", fmt.Errorf("%w: invalid password: %v", web.ErrUserInvalidCredentials, err)
	}

	return web.AuthenticationUserID(id), nil
}

func (s *SQLite) Lookup(id web.AuthenticationUserID) (web.AuthenticationUser, error) {
	rows, err := s.db.Query("SELECT id, username FROM authentication_user WHERE id = ?", id.String())
	if err != nil {
		return web.AuthenticationUser{}, fmt.Errorf("can't lookup user in database: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return web.AuthenticationUser{}, web.ErrUserNotFound
	}

	var user web.AuthenticationUser
	if err := rows.Scan(&user.ID, &user.Username); err != nil {
		return web.AuthenticationUser{}, fmt.Errorf("can't scan SQL columns: %v", err)
	}

	return user, nil
}

func (s *SQLite) Register(username string, password string) (web.AuthenticationUserID, error) {
	id := s.generateID()
	hashedPassword, salt, err := s.hashPassword(password)
	if err != nil {
		return "", fmt.Errorf("can't hash password: %v", err)
	}

	user := web.AuthenticationUser{ID: id, Username: username}

	_, err = s.db.Exec(
		"INSERT INTO authentication_user (ID, username, password, salt) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, hashedPassword, salt,
	)

	if err != nil {
		if sqliteutil.IsUniqueConstraintError(err, "authentication_user.username") {
			return "", web.ErrUserAlreadyExist
		}
		return "", fmt.Errorf("can't insert user (username=%s): %v", username, err)
	}

	return user.ID, nil
}

func (s *SQLite) hashPassword(password string) (string, string, error) {
	salt, err := s.generateSalt()
	if err != nil {
		return "", "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt+s.pepper), s.HashCost)
	if err != nil {
		return "", "", fmt.Errorf("can't generate bcrypt password: %v", err)
	}

	return string(hashedPassword), salt, nil
}

func (s *SQLite) generateSalt() (string, error) {
	buffer := make([]byte, 10)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("can't generate salt: %v", err)
	}

	return base64.URLEncoding.EncodeToString(buffer), err
}

func (s *SQLite) generateID() web.AuthenticationUserID {
	return web.AuthenticationUserID(uuid.NewString())
}

var _ web.AuthenticationBackendStorer = &SQLite{}
