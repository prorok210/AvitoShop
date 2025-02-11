package models

import (
	"time"
)

type Model struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Balance  int    `json:"balance,omitempty"`
}

type Tokens struct {
	Model
	UserID       int
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
