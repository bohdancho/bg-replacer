package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore() Store {
	db, err := sql.Open("sqlite3", "./imaginaer.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(scheme)
	if err != nil {
		panic(err)
	}

	return Store{db: db}
}

var scheme = `
CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY,
  username STRING NOT NULL UNIQUE,
  password STRING NOT NULL
);

CREATE TABLE IF NOT EXISTS session (
  token STRING PRIMARY KEY,
  expires DATETIME NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user (id)
);

CREATE TABLE IF NOT EXISTS image (
  url STRING PRIMARY KEY,
  owner_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY (owner_id) REFERENCES user (id)
);`
