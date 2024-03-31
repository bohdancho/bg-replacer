package db

import (
	"database/sql"
	"imaginaer/auth"

	"github.com/mattn/go-sqlite3"
)

func (s Store) CreateUser(user auth.User) (auth.UserID, error) {
	result, err := s.db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, auth.ErrUsernameTaken
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return auth.UserID(id), nil
}

func (s Store) DeleteUser(id auth.UserID) error {
	stmt, err := s.db.Prepare("DELETE FROM user WHERE id = ?")
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
		return auth.ErrUserNotFound
	}
	return nil
}

func (s Store) UserByID(id auth.UserID) (auth.User, error) {
	var user auth.User
	row := s.db.QueryRow("SELECT * FROM user WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return auth.User{}, auth.ErrUserNotFound
		}
		return auth.User{}, err
	}
	return user, nil
}

func (s Store) UserByUsername(username string) (auth.User, error) {
	var user auth.User

	row := s.db.QueryRow("SELECT * FROM user WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, auth.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
