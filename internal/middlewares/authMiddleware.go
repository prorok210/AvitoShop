package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/utils"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
			}

			token := parts[1]
			claims, err := utils.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}

			name, ok := claims["name"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			var userID, balance int
			err = db.DBConn.QueryRow(context.Background(),
				"SELECT user_id, balance FROM users WHERE name = $1", name).
				Scan(&userID, &balance)
			if errors.Is(err, pgx.ErrNoRows) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
			} else if err != nil {

				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
			}

			c.Set("userID", userID)
			c.Set("Name", name)
			c.Set("Balance", balance)

			return next(c)
		}
	}
}
