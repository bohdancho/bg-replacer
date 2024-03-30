package auth

import (
	"errors"
	"fmt"
	"imaginaer/db"
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
		fmt.Println(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return SessionID(id), nil
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
