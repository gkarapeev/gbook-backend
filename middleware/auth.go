package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	m "this_project_id_285410/models"
)

func AuthMiddleware(next http.Handler, db *sql.DB) http.Handler {
	unprotected := map[string]bool{
		"/register":   true,
		"/login":      true,
		"/login-auto": true,
		"/logout":     true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AUTH: requesting ", r.URL.Path)

		if unprotected[r.URL.Path] {
			log.Println("AUTH: allowing exception ", r.URL.Path)

			next.ServeHTTP(w, r)
			return
		}

		godotenv.Load()
		jSecret := os.Getenv("JWT_SECRET")

		if jSecret == "" {
			http.Error(w, "JWT secret not set", http.StatusInternalServerError)
			return
		}
		jwtKey := []byte(jSecret)

		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != "HS256" {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))

		var user m.DbUser
		err = db.QueryRow(
			"SELECT id, username FROM users WHERE id = $1",
			userID,
		).Scan(&user.ID, &user.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
