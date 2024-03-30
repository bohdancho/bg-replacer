package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "db/imaginaer.db")
	DB = db
	if err != nil {
		panic(err)
	}
}
