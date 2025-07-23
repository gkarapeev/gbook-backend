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

	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(w, r, db)
	})

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(w, r, db)
	})

	mux.HandleFunc("/api/registry", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetRegistry(w, r, db)
	})

	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPostsByUser(w, r, db)
	})

	mux.HandleFunc("/api/createPost", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, db)
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", middleware.CorsMiddleware(middleware.AuthMiddleware(mux)))
}
