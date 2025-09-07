package friendship

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"this_project_id_285410/models"
)

func CancelFriendRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	var payload models.CancelFriendRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.OtherUserID == 0 {
		http.Error(w, "Missing otherUserId", http.StatusBadRequest)
		return
	}

	// Use DB function to cancel friendship
	var result string
	// pass current user and other user; function handles canonical ordering
	err := db.QueryRow("SELECT cancel_friend_request($1, $2)", user.ID, payload.OtherUserID).Scan(&result)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": result})
}
