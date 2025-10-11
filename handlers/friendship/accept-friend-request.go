package friendship

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"this_project_id_285410/models"
)

func AcceptFriendRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(*models.DbUser)
	if user == nil {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		SenderID int `json:"senderId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.SenderID == 0 {
		http.Error(w, "Missing senderId", http.StatusBadRequest)
		return
	}

	var result string
	err := db.QueryRow("SELECT accept_friend_request($1, $2)", user.ID, payload.SenderID).Scan(&result)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": result})
}
