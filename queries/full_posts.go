package queries

import (
	"database/sql"
	"fmt"
	"strings"

	m "this_project_id_285410/models"
)

func QueryFullPosts(db *sql.DB, hostID *int, skip int, take int, likeAuthority int) ([]m.FullPost, error) {
	postIDs, err := getPostIDs(db, hostID, skip, take)
	if err != nil {
		return nil, err
	}
	if len(postIDs) == 0 {
		return []m.FullPost{}, nil
	}

	postsMap := make(map[int]*m.FullPost, len(postIDs))
	var posts []*m.FullPost

	likedPosts, err := getLikedPosts(db, postIDs, likeAuthority)
	if err != nil {
		return nil, err
	}

	err = getPosts(db, postIDs, postsMap, &posts, likedPosts)
	if err != nil {
		return nil, err
	}

	err = getPostComments(db, postIDs, postsMap)
	if err != nil {
		return nil, err
	}

	finalPosts := make([]m.FullPost, len(posts))
	for i, post := range posts {
		finalPosts[i] = *post
	}

	return finalPosts, nil
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

func getPosts(db *sql.DB, postIDs []int, postsMap map[int]*m.FullPost, postsList *[]*m.FullPost, likedPosts map[int]bool) error {
	query := `
        SELECT
            p.id, p.author_id, p.content, p.created_at, p.image_present,
            pa.username,
            p.host_id, hu.username
        FROM posts p
        JOIN users pa ON p.author_id = pa.id
        JOIN users hu ON p.host_id = hu.id
        WHERE p.id IN (%s)
        ORDER BY p.id DESC
    `
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query = fmt.Sprintf(query, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("error executing posts query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postID, postAuthorID, postHostID sql.NullInt64
		var postContent, postAuthorUsername, postHostUsername sql.NullString
		var postCreatedAt sql.NullTime
		var postImagePresent sql.NullBool

		if err := rows.Scan(
			&postID, &postAuthorID, &postContent, &postCreatedAt, &postImagePresent,
			&postAuthorUsername,
			&postHostID, &postHostUsername,
		); err != nil {
			return fmt.Errorf("error scanning post row: %w", err)
		}

		newPost := &m.FullPost{
			ID:           int(postID.Int64),
			Content:      postContent.String,
			CreatedAt:    int(postCreatedAt.Time.Unix() * 1000),
			ImagePresent: postImagePresent.Bool,
			Author: m.DbUser{
				ID:       int(postAuthorID.Int64),
				BaseUser: m.BaseUser{Username: postAuthorUsername.String},
			},
			Host: m.DbUser{
				ID:       int(postHostID.Int64),
				BaseUser: m.BaseUser{Username: postHostUsername.String},
			},
			Comments:    []m.FullComment{},
			UserLikesIt: likedPosts[int(postID.Int64)],
		}
		postsMap[int(postID.Int64)] = newPost
		*postsList = append(*postsList, newPost)
	}

	return rows.Err()
}

func getPostComments(db *sql.DB, postIDs []int, postsMap map[int]*m.FullPost) error {
	query := `
        SELECT
            c.post_id, c.id, c.author_id, ca.username, c.content, c.created_at
        FROM post_comments c
        JOIN users ca ON c.author_id = ca.id
        WHERE c.post_id IN (%s)
        ORDER BY c.created_at
    `

	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))

	for i, id := range postIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query = fmt.Sprintf(query, strings.Join(placeholders, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("error executing comments query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var commentPostID, commentID, commentAuthorID sql.NullInt64
		var commentAuthorUsername, commentContent sql.NullString
		var commentCreatedAt sql.NullTime

		if err := rows.Scan(
			&commentPostID, &commentID, &commentAuthorID, &commentAuthorUsername, &commentContent, &commentCreatedAt,
		); err != nil {
			return fmt.Errorf("error scanning comment row: %w", err)
		}

		if post, ok := postsMap[int(commentPostID.Int64)]; ok {
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

	return rows.Err()
}

func getLikedPosts(db *sql.DB, postIDs []int, likeAuthority int) (map[int]bool, error) {
	likedPosts := make(map[int]bool)
	query := `
        SELECT post_id
        FROM post_likes
        WHERE user_id = $1 AND post_id IN (%s)
    `
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs)+1)
	args[0] = likeAuthority
	for i, id := range postIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}

	query = fmt.Sprintf(query, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing liked posts query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return nil, fmt.Errorf("error scanning liked post id: %w", err)
		}
		likedPosts[postID] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during liked posts row iteration: %w", err)
	}

	return likedPosts, nil
}
