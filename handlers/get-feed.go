package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	m "this_project_id_285410/models"
	"this_project_id_285410/queries"
)

func GetFeed(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	skip := 0
	take := 50

	skipStr := r.URL.Query().Get("skip")
	takeStr := r.URL.Query().Get("take")

	if skipStr != "" {
		if val, err := strconv.Atoi(skipStr); err == nil && val >= 0 {
			skip = val
		}
	}

	if takeStr != "" {
		if val, err := strconv.Atoi(takeStr); err == nil && val > 0 {
			take = val
		}
	}

	user := r.Context().Value("user").(*m.DbUser)

	posts, err := queries.QueryFullPosts(db, nil, skip, take, user.ID)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
