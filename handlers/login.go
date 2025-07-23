package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"time"

	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	. "this_project_id_285410/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	var creds NewUser

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedHash string
	var user DbUser

	err := db.QueryRow(
		"SELECT id, username, passwordHash FROM users WHERE username = ?",
		creds.Username,
	).Scan(&user.ID, &user.Username, &storedHash)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	godotenv.Load()
	jSecret := os.Getenv("JWT_SECRET")
	if jSecret == "" {
		http.Error(w, "JWT secret not set", http.StatusInternalServerError)
		return
	}
	var jwtKey = []byte(jSecret)

	expTime := time.Now().Add(time.Minute * 15)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expTime,
		HttpOnly: true,
		Secure:   false, // set to true in production with HTTPS
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    user,
		"expires": expTime.Unix(),
	})
}
