package db

import (
	"database/sql"
	"imaginaer/auth"
	"time"
)

func (s Store) CreateSession(userID auth.UserID, expires time.Time) (auth.SessionID, error) {
	result, err := s.db.Exec("INSERT INTO session (expires, user_id) VALUES (?, ?)", expires, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return auth.SessionID(id), nil
}

func (s Store) UserBySessionId(id auth.SessionID) (auth.User, error) {
	var user auth.User
	row := s.db.QueryRow(` SELECT user.id, username, password
							from user
							INNER JOIN session ON user.id = session.user_id		
							WHERE session.id = ?;`, id)
	err := row.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return auth.User{}, auth.ErrUserNotFound
		}
		return auth.User{}, err
	}
	return user, nil
}

func (s Store) SessionByID(id auth.SessionID) (auth.Session, error) {
	var session auth.Session
	row := s.db.QueryRow("SELECT * FROM session WHERE id = ?", id)
	err := row.Scan(&session.ID, &session.Expires, &session.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return auth.Session{}, auth.ErrSessionNotFound
		}
		return auth.Session{}, err
	}
	return session, nil
}

func (s Store) DeleteSession(id auth.SessionID) error {
	stmt, err := s.db.Prepare("DELETE FROM session WHERE id = ?")
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
		return auth.ErrSessionNotFound
	}
	return nil
}
