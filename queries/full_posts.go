package queries

import (
	"database/sql"
	"fmt"
	"strings"

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
		idRows, err = db.Query(`SELECT id FROM posts WHERE host_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`, *hostID, take, skip)
	} else {
		idRows, err = db.Query(`SELECT id FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2`, take, skip)
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
	args := []interface{}{}
	placeholders := []string{}

	if hostID != nil {
		args = append(args, *hostID)
	}

	offset := 1
	if hostID != nil {
		offset = 2
	}

	for i, id := range postIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+offset)) // <-- so ridiculous lol. Sorry I can't write go yet 🐿
		args = append(args, id)
	}

	var whereClause string
	if hostID != nil {
		whereClause = fmt.Sprintf("p.host_id = $1 AND p.id IN (%s)", strings.Join(placeholders, ","))
	} else {
		whereClause = fmt.Sprintf("p.id IN (%s)", strings.Join(placeholders, ","))
	}

	query := fmt.Sprintf(`
		SELECT
			p.id, p.author_id, p.content, p.created_at, p.image_present,
			pa.id, pa.username,
			p.host_id, hu.username,
			c.id, c.post_id, c.author_id, c.content, c.created_at,
			ca.username
		FROM posts p
		JOIN users pa ON p.author_id = pa.id
		JOIN users hu ON p.host_id = hu.id
		LEFT JOIN post_comments c ON p.id = c.post_id
		LEFT JOIN users ca ON c.author_id = ca.id
		WHERE %s
		ORDER BY p.id DESC, c.created_at ASC
	`, whereClause)

	return query, args
}

func scanFullPosts(rows *sql.Rows) ([]m.FullPost, error) {
	postsMap := make(map[int]*m.FullPost)
	var posts []*m.FullPost

	for rows.Next() {
		var postID, postAuthorID, postHostID, commentID, commentPostID, commentAuthorID sql.NullInt64
		var postContent, postAuthorUsername, postHostUsername, commentContent, commentAuthorUsername sql.NullString
		var postCreatedAt, commentCreatedAt sql.NullTime
		var postImagePresent sql.NullBool

		if err := rows.Scan(
			&postID, &postAuthorID, &postContent, &postCreatedAt, &postImagePresent,
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
				CreatedAt: int(postCreatedAt.Time.Unix() * 1000),
				Author: m.DbUser{
					ID:       int(postAuthorID.Int64),
					BaseUser: m.BaseUser{Username: postAuthorUsername.String},
				},
				Host: m.DbUser{
					ID:       int(postHostID.Int64),
					BaseUser: m.BaseUser{Username: postHostUsername.String},
				},
				ImagePresent: postImagePresent.Bool,
				Comments:     []m.FullComment{},
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
				CreatedAt: commentCreatedAt.Time.Unix() * 1000,
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
