package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	. "this_project_id_285410/models"
)

func GetPostsByUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.URL.Query().Get("userId")

	if userIDStr == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, userId, content FROM posts WHERE userId = ? ORDER BY id DESC", userID)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		if err := rows.Scan(&post.ID, &post.UserID, &post.Content); err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
