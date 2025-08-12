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

func fetchFullPosts(db *sql.DB, hostID int, skip int, take int) ([]m.FullPost, error) {
	var hostUser m.DbUser
	hostRow := db.QueryRow("SELECT id, username FROM users WHERE id = ?", hostID)
	if err := hostRow.Scan(&hostUser.ID, &hostUser.Username); err != nil {
		return nil, fmt.Errorf("error fetching host user: %w", err)
	}

	// Step 1: Get post IDs for pagination
	idRows, err := db.Query(`SELECT id FROM posts WHERE hostId = ? ORDER BY id DESC LIMIT ? OFFSET ?`, hostID, take, skip)
	if err != nil {
		return nil, fmt.Errorf("error fetching post ids: %w", err)
	}
	defer idRows.Close()
	var postIDs []int
	for idRows.Next() {
		var id int
		if err := idRows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning post id: %w", err)
		}
		postIDs = append(postIDs, id)
	}
	if len(postIDs) == 0 {
		return []m.FullPost{}, nil
	}

	// Step 2: Get all post/comment data for those post IDs
	// Build placeholders for IN clause
	placeholders := ""
	args := make([]interface{}, len(postIDs)+1)
	args[0] = hostID

	for i := range postIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i+1] = postIDs[i]
	}

	query := fmt.Sprintf(`
		SELECT
			p.id, p.authorId, p.content, p.createdAt,
			pa.id, pa.username,
			c.id, c.postId, c.authorId, c.content, c.createdAt,
			ca.username
		FROM posts p
		JOIN users pa ON p.authorId = pa.id
		LEFT JOIN post_comments c ON p.id = c.postId
		LEFT JOIN users ca ON c.authorId = ca.id
		WHERE p.hostId = ? AND p.id IN (%s)
		ORDER BY p.id DESC, c.createdAt DESC
	`, placeholders)

	rows, err := db.Query(query, args...)
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
	skipStr := r.URL.Query().Get("skip")
	takeStr := r.URL.Query().Get("take")

	if userIDStr == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	skip := 0
	if skipStr != "" {
		if val, err := strconv.Atoi(skipStr); err == nil && val >= 0 {
			skip = val
		}
	}

	take := 20
	if takeStr != "" {
		if val, err := strconv.Atoi(takeStr); err == nil && val > 0 {
			take = val
		}
	}

	posts, err := fetchFullPosts(db, userID, skip, take)
	if err != nil {
		log.Printf("Error in GetPostsByUser: %v", err)
		http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
