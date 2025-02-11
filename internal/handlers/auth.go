package handlers

import (
	"context"
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
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var dbPassword string
	var userID int
	isNewUser := false

	err := db.DBConn.QueryRow(context.Background(),
		"SELECT user_id, password FROM users WHERE email = $1", u.Email).
		Scan(&userID, &dbPassword)
	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		hashedPass, err := utils.HashPassword(u.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}
		_, err = db.DBConn.Exec(context.Background(),
			"INSERT INTO users(name, email, password) VALUES ($1, $2, $3)",
			u.Name, u.Email, hashedPass)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		isNewUser = true
		dbPassword = hashedPass

		err = db.DBConn.QueryRow(context.Background(),
			"SELECT user_id FROM users WHERE email = $1", u.Email).Scan(&userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	if userID == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !utils.CheckPassword(dbPassword, u.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
	}

	access, err := utils.GenerateAccessToken(u.Name, u.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}
	refresh, err := utils.GenerateRefreshToken(u.Name, u.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	token := &models.Tokens{AccessToken: access, RefreshToken: refresh, UserID: userID}
	_, err = db.DBConn.Exec(context.Background(), "INSERT INTO tokens(access_token, refresh_token, user_id) VALUES ($1, $2, $3)",
		token.AccessToken, token.RefreshToken, token.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to store token"})
	}
	responseData := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	if isNewUser {
		return c.JSON(http.StatusCreated, responseData)
	}
	return c.JSON(http.StatusOK, responseData)
}
