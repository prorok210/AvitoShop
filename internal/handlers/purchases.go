package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
)

func Buy(c echo.Context) error {
	item := c.Param("item")

	var merch_id, price int
	err := db.DBConn.QueryRow(context.Background(),
		"SELECT merch_id, price FROM merch WHERE name = $1", item).Scan(&merch_id, &price)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	balance := c.Get("Balance").(int)
	if balance < price {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Not enough money"})
	}

	userID := c.Get("userID").(int)
	tx, err := db.DBConn.Begin(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE user_id = $2", price, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update balance"})
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO orders (user_id, merch_id) VALUES ($1, $2)", userID, merch_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert order"})
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Purchase successful"})
}
