package db

import (
	"database/sql"
	"imaginaer/auth"
	"time"
)

func (s Store) CreateSession(token auth.SessionToken, userID auth.UserID, expires time.Time) error {
	_, err := s.db.Exec("INSERT INTO session (token, expires, user_id) VALUES (?, ?, ?);", token, expires, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) UserBySessionToken(token auth.SessionToken) (auth.User, error) {
	var user auth.User
	row := s.db.QueryRow(` SELECT user.id, username, password
							from user
							INNER JOIN session ON user.id = session.user_id		
							WHERE session.token = ?;`, token)
	err := row.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return auth.User{}, auth.ErrUserNotFound
		}
		return auth.User{}, err
	}
	return user, nil
}

func (s Store) SessionByToken(token auth.SessionToken) (auth.Session, error) {
	var session auth.Session
	row := s.db.QueryRow("SELECT * FROM session WHERE token = ?", token)
	err := row.Scan(&session.Token, &session.Expires, &session.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return auth.Session{}, auth.ErrSessionNotFound
		}
		return auth.Session{}, err
	}
	return session, nil
}

func (s Store) DeleteSession(token auth.SessionToken) error {
	stmt, err := s.db.Prepare("DELETE FROM session WHERE token = ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(token)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return auth.ErrSessionNotFound
	}
	return nil
}
