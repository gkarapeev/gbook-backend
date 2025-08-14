package queries

import (
	"database/sql"
	"fmt"

	m "this_project_id_285410/models"
)

func QueryFullPosts(db *sql.DB, hostID *int, skip int, take int) ([]m.FullPost, error) {
	postIDs, err := getPostIDs(db, hostID, skip, take)

	if err != nil {
		return nil, err
	}

	if len(postIDs) == 0 {
		return []m.FullPost{}, nil
	}

	query, args := buildFullPostsQuery(hostID, postIDs)
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	return scanFullPosts(rows)
}

func getPostIDs(db *sql.DB, hostID *int, skip int, take int) ([]int, error) {
	var (
		idRows *sql.Rows
		err    error
	)

	if hostID != nil {
		idRows, err = db.Query(`SELECT id FROM posts WHERE hostId = ? ORDER BY id DESC LIMIT ? OFFSET ?`, *hostID, take, skip)
	} else {
		idRows, err = db.Query(`SELECT id FROM posts ORDER BY id DESC LIMIT ? OFFSET ?`, take, skip)
	}

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

	return postIDs, nil
}

func buildFullPostsQuery(hostID *int, postIDs []int) (string, []interface{}) {
	placeholders := ""
	args := make([]interface{}, 0, len(postIDs)+1)

	for i := range postIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, postIDs[i])
	}

	var whereClause string
	if hostID != nil {
		args = append([]interface{}{*hostID}, args...)
		whereClause = fmt.Sprintf("p.hostId = ? AND p.id IN (%s)", placeholders)
	} else {
		whereClause = fmt.Sprintf("p.id IN (%s)", placeholders)
	}

	query := fmt.Sprintf(`
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
		WHERE %s
		ORDER BY p.id DESC, c.createdAt DESC
	`, whereClause)

	return query, args
}

func scanFullPosts(rows *sql.Rows) ([]m.FullPost, error) {
	postsMap := make(map[int]*m.FullPost)
	var posts []*m.FullPost

	for rows.Next() {
		var postID, postAuthorID, postHostID, commentID, commentPostID, commentAuthorID sql.NullInt64
		var postContent, postAuthorUsername, postHostUsername, commentContent, commentAuthorUsername sql.NullString
		var postCreatedAt, commentCreatedAt sql.NullInt64

		if err := rows.Scan(
			&postID, &postAuthorID, &postContent, &postCreatedAt,
			&postAuthorID, &postAuthorUsername,
			&postHostID, &postHostUsername,
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	finalPosts := make([]m.FullPost, len(posts))
	for i, post := range posts {
		finalPosts[i] = *post
	}

	return finalPosts, nil
}
