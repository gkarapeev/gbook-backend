package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	m "this_project_id_285410/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	var user m.NewUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var existingUserID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", user.Username).Scan(&existingUserID)
	if err != sql.ErrNoRows {
		if err == nil {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	hashedPassword := string(hash)

	_, err = db.Exec("INSERT INTO users (username, passwordHash) VALUES (?, ?)", user.Username, hashedPassword)

	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}
