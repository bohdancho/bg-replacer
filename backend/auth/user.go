package auth

import (
	"database/sql"
	"errors"
	"imaginaer/db"

	"github.com/mattn/go-sqlite3"
)

type UserID int64
type User struct {
	ID       UserID
	username string
	password string
}

var ErrUsernameTaken = errors.New("username taken")

func addUser(user User) (UserID, error) {
	result, err := db.DB.Exec("INSERT INTO user (username, password) VALUES (?, ?)", user.username, user.password)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, ErrUsernameTaken
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return UserID(id), nil
}

func deleteUser(id UserID) error {
	stmt, err := db.DB.Prepare("DELETE FROM user WHERE id = ?")
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
		return ErrUserNotFound
	}
	return nil
}

var ErrUserNotFound = errors.New("user not found")

func userByUsername(username string) (User, error) {
	var user User

	row := db.DB.QueryRow("SELECT * FROM user WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.username, &user.password); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
