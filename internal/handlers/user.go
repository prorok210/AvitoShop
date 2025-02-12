package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/models"
	"github.com/prorok210/AvitoShop/internal/utils"
)

// Auth godoc
// @Summary Аутентификация и получение JWT-токена.
// @Description Аутентификация с помощью имени пользователя и пароля и возвращение access и refresh токенов.
// @Tags User
// @Accept application/json
// @Produce application/json
// @Param body body models.AuthRequest false "Auth credentials"
// @Success 200 {object} models.AuthResponse "Успешная аутентификация"
// @Failure 400 {object} models.Error400Response "Неверный запрос"
// @Failure 401 {object} models.Error401Response "Неавторизован"
// @Failure 500 {object} models.Error500Response "Внутренняя ошибка сервера"
// @Router /auth [post]
func Auth(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if u.Name == "" || u.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var dbPassword string
	var userID int
	isNewUser := false

	err := db.DBConn.QueryRow(context.Background(),
		"SELECT user_id, password FROM users WHERE name = $1", u.Name).
		Scan(&userID, &dbPassword)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			hashedPass, err := utils.HashPassword(u.Password)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
			}
			_, err = db.DBConn.Exec(context.Background(),
				"INSERT INTO users(name, password) VALUES ($1, $2)",
				u.Name, hashedPass)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
			}

			isNewUser = true
			dbPassword = hashedPass

			err = db.DBConn.QueryRow(context.Background(),
				"SELECT user_id FROM users WHERE name = $1", u.Name).Scan(&userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		} else {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	if userID == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user ID"})
	}
	if !utils.CheckPassword(dbPassword, u.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
	}

	access, err := utils.GenerateAccessToken(u.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}
	refresh, err := utils.GenerateRefreshToken(u.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	_, err = db.DBConn.Exec(context.Background(), "INSERT INTO tokens(access_token, refresh_token, user_id) VALUES ($1, $2, $3)",
		access, refresh, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to store token"})
	}
	responseData := map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	}

	if isNewUser {
		return c.JSON(http.StatusCreated, responseData)
	}
	return c.JSON(http.StatusOK, responseData)
}

// GetInfo godoc
// @Summary Получить информацию о монетах, инвентаре и истории транзакций.
// @Description Получение баланса, инвентаря и истории транзакций (отправленных и полученных монет) для авторизованного пользователя.
// @Tags User
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} models.InfoResponse "Успешный ответ"
// @Failure 400 {object} models.Error400Response "Неверный запрос"
// @Failure 401 {object} models.Error401Response "Неавторизован"
// @Failure 500 {object} models.Error500Response "Внутренняя ошибка сервера"
// @Router /info [get]
func GetInfo(c echo.Context) error {
	userID := c.Get("userID").(int)
	balance := c.Get("Balance").(int)

	inventory := []models.InventoryItem{}
	rows, err := db.DBConn.Query(context.Background(),
		"SELECT m.name, COUNT(*) FROM orders o JOIN merch m ON o.merch_id = m.merch_id WHERE o.user_id = $1 GROUP BY m.name", userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query inventory"})
	}
	defer rows.Close()
	for rows.Next() {
		var item models.InventoryItem
		var count int
		if err := rows.Scan(&item.Type, &count); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan inventory"})
		}
		item.Quantity = count
		inventory = append(inventory, item)
	}

	sent := []models.SentTx{}
	sentRows, err := db.DBConn.Query(context.Background(),
		"SELECT u.name, t.amount_coins FROM transactions t JOIN users u ON t.recipient_id = u.user_id WHERE t.user_id = $1", userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query sent transactions"})
	}
	defer sentRows.Close()
	for sentRows.Next() {
		var tx models.SentTx
		if err := sentRows.Scan(&tx.ToUser, &tx.Amount); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan sent transactions"})
		}
		sent = append(sent, tx)
	}

	received := []models.ReceivedTx{}
	recRows, err := db.DBConn.Query(context.Background(),
		"SELECT u.name, t.amount_coins FROM transactions t JOIN users u ON t.user_id = u.user_id WHERE t.recipient_id = $1", userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query received transactions"})
	}
	defer recRows.Close()
	for recRows.Next() {
		var tx models.ReceivedTx
		if err := recRows.Scan(&tx.FromUser, &tx.Amount); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan received transactions"})
		}
		received = append(received, tx)
	}

	response := models.InfoResponse{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: models.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	return c.JSON(http.StatusOK, response)
}
