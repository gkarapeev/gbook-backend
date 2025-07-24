package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	unprotected := map[string]bool{
		"/register": true,
		"/login":    true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Checking authentication for", r.URL.Path)

		if unprotected[r.URL.Path] {
			log.Println("Allowing", r.URL.Path)

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

		next.ServeHTTP(w, r)
	})
}
