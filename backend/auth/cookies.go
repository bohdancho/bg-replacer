package auth

import (
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
