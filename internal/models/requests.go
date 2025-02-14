package models

type AuthRequest struct {
	Username string `json:"username" example:"username" binding:"required"`
	Password string `json:"password" example:"secret123" binding:"required"`
}

type TransactionRequest struct {
	ToUser string `json:"toUser" example:"username" binding:"required"`
	Amount int    `json:"amount" example:"100" binding:"required"`
}
