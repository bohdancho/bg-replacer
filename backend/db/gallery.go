package db

import (
	"database/sql"
	"imaginaer/auth"
	"imaginaer/gallery"

	"github.com/mattn/go-sqlite3"
)

func (s Store) CreateImageUrl(url string, ownerID auth.UserID) error {
	_, err := s.db.Exec("INSERT INTO image (url, owner_id) VALUES (?, ?);", url, ownerID)

	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return gallery.ErrImageAlreadyUploaded
			}
		}
		return err
	}

	return nil
}

func (s Store) DeleteImageUrl(url string) error {
	stmt, err := s.db.Prepare("DELETE FROM image WHERE url = ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(url)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return gallery.ErrImageNotFound
	}
	return nil
}

func (s Store) ImageUrlByOwner(ownerId auth.UserID) (string, error) {
	var url string
	row := s.db.QueryRow("SELECT * FROM image WHERE owner_id = ?", ownerId)
	err := row.Scan(&url)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", gallery.ErrImageNotFound
		}
		return "", err
	}
	return url, nil
}
