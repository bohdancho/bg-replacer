package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore() Store {
	db, err := sql.Open("sqlite3", "db/imaginaer.db")
	if err != nil {
		panic(err)
	}
	return Store{db: db}
}
