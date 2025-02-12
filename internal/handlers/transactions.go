package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/models"
)

// SendCoin godoc
// @Summary Отправить монеты другому пользователю.
// @Description Перевод монет от одного пользователя к другому.
// @Tags Transactions
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param body body models.TransactionRequest true "SendCoinRequest"
// @Success 200 {object} models.SuccessResponse "Успешный запрос"
// @Failure 400 {object} models.Error400Response "Неверный запрос"
// @Failure 401 {object} models.Error401Response "Неавторизован"
// @Failure 500 {object} models.Error500Response "Внутренняя ошибка сервера"
// @Router /sendCoin [post]
func SendCoin(c echo.Context) error {
	req := new(models.Transaction)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ToUser == "" || req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	fromUser := c.Get("userID").(int)
	balance := c.Get("Balance").(int)
	if balance < req.Amount {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Not enough money"})
	}

	var toUserID int
	err := db.DBConn.QueryRow(context.Background(),
		"SELECT user_id FROM users WHERE name = $1", req.ToUser).Scan(&toUserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Recipient not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if toUserID == fromUser {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "You can't send coins to yourself"})
	}

	tx, err := db.DBConn.Begin(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE user_id = $2",
		req.Amount, fromUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update sender's balance"})
	}

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance + $1 WHERE user_id = $2",
		req.Amount, toUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update recipient's balance"})
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO transactions (user_id, recipient_id, amount_coins) VALUES ($1, $2, $3)",
		fromUser, toUserID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert transaction record"})
	}

	if err = tx.Commit(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Coins transferred successfully"})
}
