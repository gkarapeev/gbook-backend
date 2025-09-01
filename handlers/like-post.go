package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	m "this_project_id_285410/models"
)

func LikePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(*m.DbUser)

	var like m.FrontendLike
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if like.Unlike {
		_, err := db.Exec("DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2", like.PostID, user.ID)
		if err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "unliked"})
		return
	}

	db.QueryRow(
		"INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) RETURNING id, created_at",
		like.PostID, user.ID,
	)

	json.NewEncoder(w).Encode(map[string]string{"status": "liked"})
}
