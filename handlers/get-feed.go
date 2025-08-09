package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	m "this_project_id_285410/models"
)

func GetFeed(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, authorId, hostId, content, createdAt FROM posts ORDER BY id DESC LIMIT 50")

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var posts []m.FullPost

	for rows.Next() {
		var postID, hostID, authorID, createdAt int
		var content string

		if err := rows.Scan(&postID, &hostID, &authorID, &content, &createdAt); err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var host m.DbUser

		hostRow := db.QueryRow("SELECT id, userName FROM users WHERE id = ?", hostID)

		if err := hostRow.Scan(&host.ID, &host.Username); err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var author m.DbUser

		authorRow := db.QueryRow("SELECT id, userName FROM users WHERE id = ?", authorID)

		if err := authorRow.Scan(&author.ID, &author.Username); err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		postWithAuthor := m.FullPost{
			ID:        postID,
			Content:   content,
			CreatedAt: createdAt,
			Author:    author,
			Host:      host,
		}
		posts = append(posts, postWithAuthor)
	}

	json.NewEncoder(w).Encode(posts)
}
