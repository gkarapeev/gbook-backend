package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "this_project_id_285410/models"
	"this_project_id_285410/queries"
)

// fetchFullPostsByUser gets post IDs for a user and returns full posts using FetchFullPosts
func fetchFullPostsByUser(db *sql.DB, hostID int, skip int, take int) ([]m.FullPost, error) {
	return queries.QueryFullPosts(db, &hostID, skip, take)
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.URL.Query().Get("userId")
	skipStr := r.URL.Query().Get("skip")
	takeStr := r.URL.Query().Get("take")

	if userIDStr == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	skip := 0
	if skipStr != "" {
		if val, err := strconv.Atoi(skipStr); err == nil && val >= 0 {
			skip = val
		}
	}

	take := 20
	if takeStr != "" {
		if val, err := strconv.Atoi(takeStr); err == nil && val > 0 {
			take = val
		}
	}

	posts, err := fetchFullPostsByUser(db, userID, skip, take)
	if err != nil {
		log.Printf("Error in GetPostsByUser: %v", err)
		http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
