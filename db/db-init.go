package db

import (
	"database/sql"
	"log"

	"os"

	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	godotenv.Load()
	location := os.Getenv("DB_LOCATION")

	db, err := sql.Open("sqlite3", location)
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
			"hostId"	INTEGER NOT NULL,
			"authorId"	INTEGER NOT NULL,
			"content"	TEXT NOT NULL,
			PRIMARY KEY("id")
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
