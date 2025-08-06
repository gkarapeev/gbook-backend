package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "this_project_id_285410/models"
)

func fetchFullPosts(db *sql.DB, hostID int) ([]m.FullPost, error) {
	var hostUser m.DbUser

	hostRow := db.QueryRow("SELECT id, username FROM users WHERE id = ?", hostID)

	if err := hostRow.Scan(&hostUser.ID, &hostUser.Username); err != nil {
		return nil, fmt.Errorf("error fetching host user: %w", err)
	}

	rows, err := db.Query(`
		SELECT
			p.id, p.authorId, p.content, p.createdAt,
			pa.id, pa.username,
			c.id, c.postId, c.authorId, c.content, c.createdAt,
			ca.username
		FROM posts p
		JOIN users pa ON p.authorId = pa.id
		LEFT JOIN post_comments c ON p.id = c.postId
		LEFT JOIN users ca ON c.authorId = ca.id
		WHERE p.hostId = ?
		ORDER BY p.id DESC, c.createdAt DESC
	`, hostID)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	var posts []*m.FullPost
	postsMap := make(map[int]*m.FullPost)

	for rows.Next() {
		var postID, postAuthorID, commentID, commentPostID, commentAuthorID sql.NullInt64
		var postContent, postAuthorUsername, commentContent, commentAuthorUsername sql.NullString
		var postCreatedAt, commentCreatedAt sql.NullInt64

		if err := rows.Scan(
			&postID, &postAuthorID, &postContent, &postCreatedAt,
			&postAuthorID, &postAuthorUsername,
			&commentID, &commentPostID, &commentAuthorID, &commentContent, &commentCreatedAt,
			&commentAuthorUsername,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		post, ok := postsMap[int(postID.Int64)]
		if !ok {
			newPost := &m.FullPost{
				ID:        int(postID.Int64),
				Content:   postContent.String,
				CreatedAt: int(postCreatedAt.Int64),
				Author: m.DbUser{
					ID:       int(postAuthorID.Int64),
					BaseUser: m.BaseUser{Username: postAuthorUsername.String},
				},
				Host:     hostUser,
				Comments: []m.FullComment{},
			}
			posts = append(posts, newPost)
			postsMap[int(postID.Int64)] = newPost
			post = newPost
		}

		if commentID.Valid {
			fullComment := m.FullComment{
				ID:        int(commentID.Int64),
				PostID:    int(commentPostID.Int64),
				Content:   commentContent.String,
				CreatedAt: commentCreatedAt.Int64,
				Author: m.DbUser{
					ID:       int(commentAuthorID.Int64),
					BaseUser: m.BaseUser{Username: commentAuthorUsername.String},
				},
			}
			post.Comments = append(post.Comments, fullComment)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	finalPosts := make([]m.FullPost, len(posts))
	for i, post := range posts {
		finalPosts[i] = *post
	}

	return finalPosts, nil
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	posts, err := fetchFullPosts(db, userID)
	if err != nil {
		log.Printf("Error in GetPostsByUser: %v", err)
		http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
