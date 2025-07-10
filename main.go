package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// test
func main() {
	db := InitDB()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		getUsersHandler(w, r, db)
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))
}
