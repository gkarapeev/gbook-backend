package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	m "this_project_id_285410/models"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	// Expect path: /users/{id}
	path := r.URL.Path
	prefix := "/users/"
	if len(path) <= len(prefix) || path[:len(prefix)] != prefix {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	id := path[len(prefix):]
	if id == "" {
		http.Error(w, "Missing user id", http.StatusBadRequest)
		return
	}

	var user m.DbUser
	err := db.QueryRow("SELECT id, username FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
