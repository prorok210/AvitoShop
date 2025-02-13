package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/models"
	"github.com/prorok210/AvitoShop/internal/utils"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, models.Error401Response{Error: "Не найден заголовок Authorization"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, models.Error401Response{Error: "Неверный формат токена"})
			}

			token := parts[1]
			claims, err := utils.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.Error401Response{Error: "Невалидный токен"})
			}

			name, ok := claims["name"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, models.Error401Response{Error: "Невалидный токен"})
			}

			var userID, balance int
			err = db.DBConn.QueryRow(context.Background(),
				"SELECT user_id, balance FROM users WHERE name = $1", name).
				Scan(&userID, &balance)
			if errors.Is(err, pgx.ErrNoRows) {
				return c.JSON(http.StatusUnauthorized, models.Error401Response{Error: "Пользователь не найден"})
			} else if err != nil {
				return c.JSON(http.StatusInternalServerError, models.Error500Response{Error: "Ошибка базы данных"})
			}

			c.Set("userID", userID)
			c.Set("Name", name)
			c.Set("Balance", balance)

			return next(c)
		}
	}
}
