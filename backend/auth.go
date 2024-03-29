package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var users = map[string]string{
	"admin": "1",
}

var sessions = map[string]Session{}

type Session struct {
	expires time.Time
	login   string
}

type registrationDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(sessions)
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

const sessionCookieName = "session_token"
const sessionMaxAge = time.Hour * 24 * 14
const sessionMaxAgeSeconds = int(sessionMaxAge) / int(time.Second)

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		login:   payload.Login,
		expires: expires,
	}
	setSessionCookie(w, sessionId)
	w.WriteHeader(http.StatusOK)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie(sessionCookieName)
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

func deleteSessionCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func setSessionCookie(w http.ResponseWriter, sessionId string) {
	sessionCookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionId,
		MaxAge:   sessionMaxAgeSeconds,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &sessionCookie)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionId := sessionCookie.Value

	session, ok := sessions[sessionId]
	hasExpired := ok && session.expires.Compare(time.Now()) == -1

	if !ok || hasExpired {
		deleteSessionCookie(w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
