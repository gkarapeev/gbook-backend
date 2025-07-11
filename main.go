package main

import (
	"log"
	"net/http"

	"this_project_id_285410/db"
	"this_project_id_285410/handlers"
	"this_project_id_285410/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := db.InitDB()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUsersHandler(w, r, db)
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", middleware.CorsMiddleware(mux))
}
