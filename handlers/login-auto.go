package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	m "this_project_id_285410/models"
)

func LoginUserAuto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Missing token cookie", http.StatusUnauthorized)
		return
	}

	godotenv.Load()
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

	var userID = claims["user_id"]
	var user m.DbUser
	err = db.QueryRow("SELECT id, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var expiryTimeInt = int64(claims["exp"].(float64))
	var expiryTimeUnix = time.Unix(expiryTimeInt, 0)

	if time.Until(expiryTimeUnix) < 24*time.Hour {
		expiryTimeUnix = time.Now().Add(time.Hour * 48)

		var newClaims = jwt.MapClaims{
			"user_id": user.ID,
			"exp":     expiryTimeUnix.Unix(),
		}

		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		tokenString, err := newToken.SignedString([]byte(jSecret))
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		secureFlag := true
		if os.Getenv("COOKIE_SECURE") == "false" {
			secureFlag = false
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Expires:  expiryTimeUnix,
			HttpOnly: true,
			Secure:   secureFlag,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    user,
		"expires": expiryTimeUnix.UnixMilli(),
	})
}
