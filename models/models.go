package models

type BaseUser struct {
	Userame string `json:"username"`
}

type NewUser struct {
	BaseUser
	Password string `json:"password"`
}

type DbUser struct {
	BaseUser
	ID           int
	PasswordHash string `json:"-"`
}
