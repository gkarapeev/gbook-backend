package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	m "this_project_id_285410/models"
)

func canViewUser(db *sql.DB, viewerID, targetID int) (bool, error) {
	if viewerID == targetID {
		return true, nil
	}

	var areFriends bool
	err := db.QueryRow(`SELECT are_friends($1, $2)`, viewerID, targetID).Scan(&areFriends)
	if err != nil {
		return false, fmt.Errorf("error checking friendship: %w", err)
	}

	return areFriends, nil
}

func UserPostsMiddleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*m.DbUser)
		if user == nil {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/posts") {
			userIDStr := r.URL.Query().Get("userId")

			if userIDStr != "" {
				userID, err := strconv.Atoi(userIDStr)
				if err != nil {
					http.Error(w, "Invalid userId", http.StatusBadRequest)
					return
				}

				canView, err := canViewUser(db, user.ID, userID)
				if err != nil {
					http.Error(w, "Error checking permissions", http.StatusInternalServerError)
					return
				}

				if !canView {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
