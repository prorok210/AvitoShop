package models

type User struct {
	ID       int
	Name     string `json:"username"`
	Password string `json:"password,omitempty"`
	Balance  int    `json:"balance,omitempty"`
}

type Tokens struct {
	ID          int
	UserID      int
	AccessToken string `json:"token"`
}

type Order struct {
	ID      int
	UserID  int
	MerchID int
}

type Transaction struct {
	ID       int
	UserID   int
	ToUserID int
	ToUser   string `json:"toUser"`
	Amount   int    `json:"amount"`
}
