package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const sessionCookieName = "session_token"
const sessionMaxAge = time.Hour * 24 * 14
const sessionMaxAgeSeconds = int(sessionMaxAge) / int(time.Second)

type cookieMaxAge int

const (
	sessionCookieMaxAge = cookieMaxAge(60 * 60 * 24 * 14)
	cookieMaxAgeDel     = cookieMaxAge(-1)
)

func sessionIDFromCookie(r *http.Request) (SessionID, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return 0, fmt.Errorf("sessionIDFromCookie: %v", err)
	}
	idInt, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, ErrInvalidSessionCookie
	}
	return SessionID(idInt), nil
}

func removeSessionCookie(w http.ResponseWriter) {
	c := NewSessionCookie(0, cookieMaxAgeDel)
	http.SetCookie(w, &c)
}

func setSessionCookie(w http.ResponseWriter, sessionID SessionID) {
	c := NewSessionCookie(sessionID, sessionCookieMaxAge)
	http.SetCookie(w, &c)
}

func NewSessionCookie(sessionID SessionID, maxAge cookieMaxAge) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    fmt.Sprint(sessionID),
		MaxAge:   int(maxAge),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
}
