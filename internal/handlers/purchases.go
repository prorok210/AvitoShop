package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/models"
)

// Buy godoc
// @Summary Купить предмет за монеты.
// @Description Покупка предмета за монеты: списывается стоимость предмета с баланса пользователя и создается заказ.
// @Tags Merch
// @Security BearerAuth
// @Produce application/json
// @Param item path string true "Название предмета"
// @Success 200 {object} models.SuccessResponse "Успешный запрос"
// @Failure 400 {object} models.Error400Response "Неверный запрос."
// @Failure 401 {object} models.Error401Response "Неавторизован."
// @Failure 404 {object} models.Error404Response "Предмет не найден."
// @Failure 500 {object} models.Error500Response "Внутренняя ошибка сервера."
// @Router /api/buy/{item} [get]
func Buy(c echo.Context) error {
	item := c.Param("item")

	var merch_id, price int
	err := db.DBConn.QueryRow(context.Background(),
		"SELECT merch_id, price FROM merch WHERE name = $1", item).Scan(&merch_id, &price)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return c.JSON(http.StatusNotFound, models.Error404Response{Error: "Предмет не найден"})
		}
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Ошибка базы данных"})
	}

	balance := c.Get("Balance").(int)
	if balance < price {
		return c.JSON(http.StatusBadRequest, models.Error400Response{Error: "Недостаточно средств"})
	}

	userID := c.Get("userID").(int)
	tx, err := db.DBConn.Begin(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Ошибка начала транзакции"})
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE user_id = $2", price, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Не удалось обновить баланс"})
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO orders (user_id, merch_id) VALUES ($1, $2)", userID, merch_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Не удалось создать заказ"})
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Не удалось завершить транзакцию"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "Покупка успешно совершена"})
}
