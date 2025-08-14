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
			username TEXT NOT NULL,
			passwordHash TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"id"	INTEGER UNIQUE,
			"hostId"	INTEGER NOT NULL,
			"authorId"	INTEGER NOT NULL,
			"content"	TEXT NOT NULL,
			"createdAt" INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			"updatedAt" INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			PRIMARY KEY("id")
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "post_comments" (
			"id"        INTEGER PRIMARY KEY,
			"postId"    INTEGER NOT NULL,
			"authorId"  INTEGER NOT NULL,
			"content"   TEXT NOT NULL,
			"createdAt" INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			FOREIGN KEY("postId") REFERENCES "posts"("id") ON DELETE CASCADE
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
