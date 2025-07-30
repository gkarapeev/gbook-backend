package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"os"

	"github.com/golang-jwt/jwt/v5"

	. "this_project_id_285410/models"
)

func LoginUserAuto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Missing token cookie", http.StatusUnauthorized)
		return
	}

	jSecret := os.Getenv("JWT_SECRET")
	if jSecret == "" {
		http.Error(w, "JWT secret not set", http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(jSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return
	}

	userID := claims["user_id"]

	var user DbUser
	err = db.QueryRow("SELECT id, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
