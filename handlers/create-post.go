package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	. "this_project_id_285410/models"
)

func CreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.Content == "" {
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO posts (authorId, hostId, content) VALUES (?, ?, ?)", post.AuthorID, post.HostID, post.Content)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	post.ID = int(id)
	json.NewEncoder(w).Encode(post)
}
