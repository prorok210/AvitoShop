package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/db/mocks"
	"github.com/prorok210/AvitoShop/internal/handlers"
	"github.com/prorok210/AvitoShop/internal/models"
)

func TestSendCoin_Success(t *testing.T) {
	e := echo.New()

	reqBody, _ := json.Marshal(models.Transaction{
		ToUser: "Alice",
		Amount: 50,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/sendcoin", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", 1)
	c.Set("Balance", 100)

	mockDB := new(mocks.DB)

	mockDB.
		On("QueryRow", context.Background(),
			"SELECT user_id FROM users WHERE name = $1", "Alice").
		Return(&FakeUserRow{UserID: 2})

	mockTx := new(Tx)
	mockDB.
		On("Begin", context.Background()).
		Return(mockTx, nil)

	mockTx.
		On("Exec", context.Background(),
			"UPDATE users SET balance = balance - $1 WHERE user_id = $2",
			50, 1).
		Return(pgconn.CommandTag{}, nil)

	mockTx.
		On("Exec", context.Background(),
			"UPDATE users SET balance = balance + $1 WHERE user_id = $2",
			50, 2).
		Return(pgconn.CommandTag{}, nil)

	mockTx.
		On("Exec", context.Background(),
			"INSERT INTO transactions (user_id, recipient_id, amount_coins) VALUES ($1, $2, $3)",
			1, 2, 50).
		Return(pgconn.CommandTag{}, nil)

	mockTx.
		On("Commit", context.Background()).
		Return(nil)

	mockTx.
		On("Rollback", context.Background()).
		Return(nil)

	db.DBConn = mockDB

	err := handlers.SendCoin(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Транзакция успешно завершена")
	}
}

func TestSendCoin_InsufficientFunds(t *testing.T) {
	e := echo.New()

	reqBody, _ := json.Marshal(models.Transaction{
		ToUser: "Alice",
		Amount: 50,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/sendcoin", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", 1)
	c.Set("Balance", 20)

	err := handlers.SendCoin(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Недостаточно средств")
	}
}
