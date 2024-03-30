package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"imaginaer/db"
	"net/http"
	"strconv"
	"time"
)

type SessionID int64
type Session struct {
	ID      SessionID
	expires time.Time
	userID  UserID
}

var ErrSessionNotFound = errors.New("session not found")

func createSession(userID UserID) (SessionID, error) {
	expires := time.Now().Add(sessionMaxAge)
	result, err := db.DB.Exec("INSERT INTO session (expires, user_id) VALUES (?, ?)", expires, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return SessionID(id), nil
}

func userBySessionId(id SessionID) (User, error) {
	var user User
	row := db.DB.QueryRow(` SELECT user.id, username, password
							from user
							INNER JOIN session ON user.id = session.user_id		
							WHERE session.id = ?;`, id)
	err := row.Scan(&user.ID, &user.username, &user.password)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func sessionByID(id SessionID) (Session, error) {
	var session Session
	row := db.DB.QueryRow("SELECT * FROM session WHERE id = ?", id)
	err := row.Scan(&session.ID, &session.expires, &session.userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}
	return session, nil
}

func deleteSession(id SessionID) error {
	stmt, err := db.DB.Prepare("DELETE FROM session WHERE id = ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

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
		return 0, fmt.Errorf("sessionIDFromCookie: %v", err)
	}
	return SessionID(idInt), nil
}

func deleteSessionCookie(w http.ResponseWriter) {
	c := newSessionCookie(0, cookieMaxAgeDel)
	http.SetCookie(w, &c)
}

func setSessionCookie(w http.ResponseWriter, sessionID SessionID) {
	c := newSessionCookie(sessionID, sessionCookieMaxAge)
	http.SetCookie(w, &c)
}

func newSessionCookie(sessionID SessionID, maxAge cookieMaxAge) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    fmt.Sprint(sessionID),
		MaxAge:   int(maxAge),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
}
