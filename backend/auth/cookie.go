package auth

import (
	"fmt"
	"net/http"
	"time"
)

const sessionCookieName = "session_token"

type cookieMaxAge int

const (
	cookieSessionMaxAge = cookieMaxAge(60 * 60 * 24 * 14)
	cookieMaxAgeDel     = cookieMaxAge(-1)
)

func sessionTokenFromCookie(r *http.Request) (SessionToken, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", fmt.Errorf("sessionTokenFromCookie: %w", err)
	}
	return SessionToken(cookie.Value), nil
}

func removeSessionCookie(w http.ResponseWriter) {
	c := NewSessionCookie("del", cookieMaxAgeDel)
	http.SetCookie(w, &c)
}

func setSessionCookie(w http.ResponseWriter, sessionToken SessionToken) {
	c := NewSessionCookie(sessionToken, cookieSessionMaxAge)
	http.SetCookie(w, &c)
}

func NewSessionCookie(sessionToken SessionToken, maxAge cookieMaxAge) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    string(sessionToken),
		MaxAge:   int(maxAge),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
}

func newSessionExpiresTime() time.Time {
	t := time.Second * time.Duration(cookieSessionMaxAge)
	return time.Now().Add(t)
}
