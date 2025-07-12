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
	ID      int    `json:"id"`
	UserID  int    `json:"userId"`
	Content string `json:"content"`
}
