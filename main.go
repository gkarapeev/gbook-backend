package main

import (
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"

	"this_project_id_285410/db"
	"this_project_id_285410/handlers"
	"this_project_id_285410/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := db.InitDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(w, r, db)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutUser(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(w, r, db)
	})

	mux.HandleFunc("/login-auto", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUserAuto(w, r, db)
	})

	mux.HandleFunc("/registry", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetRegistry(w, r, db)
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUser(w, r, db)
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPostsByUser(w, r, db)
	})

	mux.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetFeed(w, r, db)
	})

	mux.HandleFunc("/createPost", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, db)
	})

	mux.HandleFunc("/addComment", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddComment(w, r, db)
	})

	mux.HandleFunc("/upload-image", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadImageHandler(w, r)
	})

	godotenv.Load()
	port := os.Getenv("PORT")
	log.Println("Server running on: " + port)
	http.ListenAndServe(":"+port, middleware.CorsMiddleware(middleware.AuthMiddleware(mux, db)))
}
