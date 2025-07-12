package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"id"	INTEGER UNIQUE,
			"userId"	INTEGER NOT NULL,
			"content"	TEXT NOT NULL,
			PRIMARY KEY("id")
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
