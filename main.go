package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)
 
func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
  )`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM users")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("["))
		first := true
		for rows.Next() {
			var id int
			var name string
			rows.Scan(&id, &name)
			if !first {
				w.Write([]byte(","))
			}
			first = false
			w.Write([]byte(`{"id":` +
				fmt.Sprintf("%d", id) + `,"name":"` + name + `"}`))
		}
		w.Write([]byte("]"))
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
