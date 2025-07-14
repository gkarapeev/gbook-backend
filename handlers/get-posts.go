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

	rows, err := db.Query("SELECT id, authorId, content, createdAt FROM posts WHERE hostId = ? ORDER BY id DESC", userID)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var posts []PostWithAuthor

	for rows.Next() {
		var postID, authorID, createdAt int
		var content string

		if err := rows.Scan(&postID, &authorID, &content, &createdAt); err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var user DbUser

		userRow := db.QueryRow("SELECT id, userName FROM users WHERE id = ?", authorID)

		if err := userRow.Scan(&user.ID, &user.Username); err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		postWithAuthor := PostWithAuthor{
			Post: Post{
				ID:        postID,
				HostID:    userID,
				AuthorID:  authorID,
				Content:   content,
				CreatedAt: createdAt,
			},
			Author: user,
		}
		posts = append(posts, postWithAuthor)
	}

	json.NewEncoder(w).Encode(posts)
}
