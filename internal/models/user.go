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
	Password string `json:"password,omitempty"`
	Balance  int    `json:"balance,omitempty"`
}

type Tokens struct {
	Model
	UserID       int
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Order struct {
	Model
	UserID  int
	MerchID int
}

type Transaction struct {
	Model
	UserID   int
	ToUserID int
	ToUser   string `json:"toUser"`
	Amount   int    `json:"amount"`
}
