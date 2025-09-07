package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	m "this_project_id_285410/models"
)

func GetRegistry(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value("user").(*m.DbUser)
	if user == nil {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
		SELECT
			u.id,
			u.username,
			f.status,
			f.sent_by_user_id
		FROM
			users u
		LEFT JOIN
			friendships f ON (LEAST(u.id, $1) = f.user_id_1 AND GREATEST(u.id, $1) = f.user_id_2)
		WHERE
			u.id != $1
	`, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []m.DbUser

	for rows.Next() {
		var regUser m.DbUser
		var status sql.NullString
		var sentByUserId sql.NullInt64

		if err := rows.Scan(&regUser.ID, &regUser.Username, &status, &sentByUserId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if status.Valid {
			regUser.FriendshipStatus = &status.String
		}
		if sentByUserId.Valid {
			sentByID := int(sentByUserId.Int64)
			regUser.SentByUserID = &sentByID
		}

		users = append(users, regUser)
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
