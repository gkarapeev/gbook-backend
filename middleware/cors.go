package middleware

import (
	"log"
	"net/http"
)

var allowedOrigins = []string{
	"https://gbook.lol",
	"https://www.gbook.lol",

	"http://test.gbook.lol",
	"http://test.www.gbook.lol",

	"http://localhost:4200",
	"http://localhost",
	"http://localhost:80",
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedSpecial := false

		for _, o := range allowedOrigins {
			if o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				allowedSpecial = true
				break
			}
		}

		if !allowedSpecial && origin != "" {
			log.Println("CORS: Origin not allowed: ", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" { // move this to the check above?
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
