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
	return mux
}

type UserID int64
type User struct {
	ID       UserID
	Username string
	Password string // TODO: remove password from here
}

type SessionID int64
type Session struct {
	ID      SessionID
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
	CreateSession(userID UserID, expires time.Time) (SessionID, error)
	UserBySessionId(id SessionID) (User, error)
	SessionByID(id SessionID) (Session, error)
	DeleteSession(id SessionID) error
}

// TODO: put this in context?
func GetCurrentUser(w http.ResponseWriter, r *http.Request, store Store) (User, error) {
	sessionID, err := sessionIDFromCookie(r)
	if err != nil {
		if err == ErrInvalidSessionCookie {
			removeSessionCookie(w)
		}
		return User{}, err
	}

	return store.UserBySessionId(sessionID)
}
