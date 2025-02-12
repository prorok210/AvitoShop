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
