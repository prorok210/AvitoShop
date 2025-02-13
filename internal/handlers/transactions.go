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
// @Failure 404 {object} models.Error404Response "Не найдено"
// @Failure 500 {object} models.Error500Response "Внутренняя ошибка сервера"
// @Router /sendCoin [post]
func SendCoin(c echo.Context) error {
	req := new(models.Transaction)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error400Response{Error: "Неверный запрос"})
	}

	if req.ToUser == "" || req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, models.Error400Response{Error: "Неверный запрос"})
	}

	fromUser := c.Get("userID").(int)
	balance := c.Get("Balance").(int)
	if balance < req.Amount {
		return c.JSON(http.StatusBadRequest, models.Error400Response{Error: "Недостаточно средств"})
	}

	var toUserID int
	err := db.DBConn.QueryRow(context.Background(),
		"SELECT user_id FROM users WHERE name = $1", req.ToUser).Scan(&toUserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNotFound, models.Error404Response{Error: "Получатель не найден"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Ошибка запроса к базе данных"})
	}
	if toUserID == fromUser {
		return c.JSON(http.StatusBadRequest, models.Error400Response{Error: "Вы не можете отправлять монеты самому себе"})
	}

	tx, err := db.DBConn.Begin(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Ошибка начала транзакции"})
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE user_id = $2",
		req.Amount, fromUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Не удалось обновить баланс"})
	}

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance + $1 WHERE user_id = $2",
		req.Amount, toUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Не удалось обновить баланс"})
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO transactions (user_id, recipient_id, amount_coins) VALUES ($1, $2, $3)",
		fromUser, toUserID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Не удалось создать транзакцию"})
	}

	if err = tx.Commit(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось завершить транзакцию"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "Транзакция успешно завершена"})
}
