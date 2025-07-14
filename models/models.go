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

type Post struct {
	ID        int    `json:"id"`
	HostID    int    `json:"hostId"`
	AuthorID  int    `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
}

type PostWithAuthor struct {
	Post
	Author DbUser `json:"author"`
}
