package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	m "this_project_id_285410/models"
)

func AddComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var comment m.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if comment.Content == "" {
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}

	if comment.AuthorID == 0 {
		http.Error(w, "Missing authorId", http.StatusBadRequest)
		return
	}

	if comment.PostID == 0 {
		http.Error(w, "Missing postId", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO post_comments (postId, authorId, content, createdAt) VALUES (?, ?, ?, ?)",
		comment.PostID, comment.AuthorID, comment.Content, time.Now().Unix())
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "DB error getting last ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	comment.ID = int(lastID)
	comment.CreatedAt = time.Now().Unix()

	var author m.DbUser
	authorRow := db.QueryRow("SELECT id, username FROM users WHERE id = ?", comment.AuthorID)
	if err := authorRow.Scan(&author.ID, &author.BaseUser.Username); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Author with ID %d not found for comment %d", comment.AuthorID, comment.ID)
			http.Error(w, "Author not found", http.StatusNotFound)
		} else {
			http.Error(w, "DB error fetching author: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fullComment := m.FullComment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		Author:    author,
	}

	json.NewEncoder(w).Encode(fullComment)
}
