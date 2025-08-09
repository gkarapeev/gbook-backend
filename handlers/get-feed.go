package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	m "this_project_id_285410/models"
)

func GetFeed(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`
		SELECT
			p.id, p.authorId, p.content, p.createdAt,
			pa.id, pa.username,
			p.hostId, hu.username,
			c.id, c.postId, c.authorId, c.content, c.createdAt,
			ca.username
		FROM posts p
		JOIN users pa ON p.authorId = pa.id
		JOIN users hu ON p.hostId = hu.id
		LEFT JOIN post_comments c ON p.id = c.postId
		LEFT JOIN users ca ON c.authorId = ca.id
		ORDER BY p.id DESC, c.createdAt DESC
		LIMIT 50
	`)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	postsMap := make(map[int]*m.FullPost)
	var posts []*m.FullPost

	for rows.Next() {
		var postID, postAuthorID, postHostID, commentID, commentPostID, commentAuthorID sql.NullInt64
		var postContent, postAuthorUsername, postHostUsername, commentContent, commentAuthorUsername sql.NullString
		var postCreatedAt, commentCreatedAt sql.NullInt64

		if err := rows.Scan(
			&postID, &postAuthorID, &postContent, &postCreatedAt,
			&postHostID, &postHostUsername,
			&postAuthorID, &postAuthorUsername,
			&commentID, &commentPostID, &commentAuthorID, &commentContent, &commentCreatedAt,
			&commentAuthorUsername,
		); err != nil {
			http.Error(w, "DB error scanning: "+err.Error(), http.StatusInternalServerError)
			return
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
				Host: m.DbUser{
					ID:       int(postHostID.Int64),
					BaseUser: m.BaseUser{Username: postHostUsername.String},
				},
				Comments: []m.FullComment{},
			}
			postsMap[int(postID.Int64)] = newPost
			posts = append(posts, newPost)
			post = newPost
		}

		if commentID.Valid {
			post.Comments = append(post.Comments, m.FullComment{
				ID:        int(commentID.Int64),
				PostID:    int(commentPostID.Int64),
				Content:   commentContent.String,
				CreatedAt: commentCreatedAt.Int64,
				Author: m.DbUser{
					ID:       int(commentAuthorID.Int64),
					BaseUser: m.BaseUser{Username: commentAuthorUsername.String},
				},
			})
		}
	}

	finalPosts := make([]m.FullPost, len(posts))
	for i, post := range posts {
		finalPosts[i] = *post
	}

	json.NewEncoder(w).Encode(finalPosts)
}
