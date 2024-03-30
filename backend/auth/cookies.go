package auth

import (
	"fmt"
	"net/http"
	"time"
)

const sessionCookieName = "session_token"
const sessionMaxAge = time.Hour * 24 * 14
const sessionMaxAgeSeconds = int(sessionMaxAge) / int(time.Second)

type sessionCookie http.Cookie

func getSessionCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(sessionCookieName)
}

func deleteSessionCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, c)
}

func setSessionCookie(w http.ResponseWriter, sessionID SessionID) {
	sessionCookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    fmt.Sprint(sessionID),
		MaxAge:   sessionMaxAgeSeconds,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &sessionCookie)
}
