package friendship

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"this_project_id_285410/models"
)

func SendFriendRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	var payload models.FriendRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.ReceiverID == 0 {
		http.Error(w, "Missing receiverId", http.StatusBadRequest)
		return
	}

	var result string
	err := db.QueryRow("SELECT send_friend_request($1, $2)", user.ID, payload.ReceiverID).Scan(&result)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": result})
}
