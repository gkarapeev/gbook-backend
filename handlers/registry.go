package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	m "this_project_id_285410/models"
)

func GetRegistry(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []m.DbUser

	for rows.Next() {
		var user m.DbUser

		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
