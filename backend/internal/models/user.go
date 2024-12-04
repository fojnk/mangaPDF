package models

type User struct {
	Guid     string `json:"guid" db:"guid"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email"  db:"email"`
	Password string `json:"-" db:"password_hash"`
}

type Token struct {
	Id        int    `json:"id" db:"id"`
	UserId    string `json:"user_id" db:"user_id"`
	TokenHash string `josn:"token_hash" db:"token_hash"`
}
