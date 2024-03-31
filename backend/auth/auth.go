package auth

import (
	"errors"
	"net/http"
	"time"
)

func NewMux(store Store) http.Handler {
	mux := http.NewServeMux()
	server := Server{store: store}

	mux.HandleFunc("/login", server.loginHandler)
	mux.HandleFunc("/logout", server.logoutHandler)
	mux.HandleFunc("/registration", server.registrationHandler)
	mux.HandleFunc("/current_user", server.currentUserHandler)
	return mux
}

type UserID int64
type User struct {
	ID       UserID
	Username string
	Password string // TODO: remove password from here
}

type UserDTO struct {
	ID       UserID `json:"id"`
	Username string `json:"username"`
}

type SessionToken string
type Session struct {
	Token   SessionToken
	Expires time.Time
	UserID  UserID
}

var ErrUserNotFound = errors.New("user not found")
var ErrUsernameTaken = errors.New("username taken")

var ErrSessionNotFound = errors.New("session not found")
var ErrInvalidSessionCookie = errors.New("invalid session cookie")

type Server struct {
	store Store
}

type Store interface {
	UserStore
	SessionStore
}

type UserStore interface {
	CreateUser(user User) (UserID, error)
	DeleteUser(id UserID) error
	UserByID(id UserID) (User, error)
	UserByUsername(username string) (User, error)
}

type SessionStore interface {
	CreateSession(token SessionToken, userID UserID, expires time.Time) error
	UserBySessionToken(token SessionToken) (User, error)
	SessionByToken(token SessionToken) (Session, error)
	DeleteSession(token SessionToken) error
}
