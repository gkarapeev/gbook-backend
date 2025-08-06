package models

type BaseUser struct {
	Username string `json:"username"`
}

type NewUser struct {
	BaseUser
	Password string `json:"password"`
}

type DbUser struct {
	BaseUser
	ID           int    `json:"id"`
	PasswordHash string `json:"-"`
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"postId"`
	AuthorID  int    `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

type FullComment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"postId"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
	Author    DbUser `json:"author"`
}

type Post struct {
	ID        int    `json:"id"`
	HostID    int    `json:"hostId"`
	AuthorID  int    `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
}

type FullPost struct {
	ID        int           `json:"id"`
	Host      DbUser        `json:"host"`
	Author    DbUser        `json:"author"`
	Content   string        `json:"content"`
	CreatedAt int           `json:"createdAt"`
	UpdatedAt int           `json:"updatedAt"`
	Comments  []FullComment `json:"comments"`
}
