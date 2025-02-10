package models

import "time"

type Model struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Model
	Name     string
	Login    string
	Password string
	Balance  int
}

type Tokens struct {
	Model
	UserID        int
	Access_token  string
	Refresh_token string

	ExpiredAt time.Time
}
