package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	login string
}
type Session struct {
	expires time.Time
	user    User
}

var users = map[string]string{
	"admin": "1",
}
var sessions = map[string]Session{}

type registrationDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload registrationDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, loginTaken := users[payload.Login]
	if loginTaken {
		w.WriteHeader(http.StatusConflict)
		return
	}

	users[payload.Login] = payload.Password
	w.WriteHeader(http.StatusOK)
}

type loginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload loginDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[payload.Login]
	if !ok || expectedPassword != payload.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId := uuid.New().String()
	expires := time.Now().Add(sessionMaxAge)

	sessions[sessionId] = Session{
		user:    User{login: payload.Login},
		expires: expires,
	}
	setSessionCookie(w, sessionId)
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := getSessionCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionId := sessionCookie.Value
	_, ok := sessions[sessionId]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessions[sessionId] = Session{}
	deleteSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

func GetUser(r *http.Request) (User, error) {
	sessionCookie, err := getSessionCookie(r)
	if err != nil {
		return User{}, errors.New("no session cookie found")
	}
	sessionId := sessionCookie.Value

	session, ok := sessions[sessionId]

	if !ok {
		return User{}, errors.New("no session matches the session cookie")
	}

	return session.user, nil
}
