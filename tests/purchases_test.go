package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/db/mocks"
	"github.com/prorok210/AvitoShop/internal/handlers"
)

func TestBuy_Success(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/buy/cup", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item")
	c.SetParamValues("cup")
	c.Set("Balance", 150)
	c.Set("userID", 99)

	mockDB := new(mocks.DB)

	mockDB.
		On("QueryRow", context.Background(),
			"SELECT merch_id, price FROM merch WHERE name = $1", "cup").
		Return(&FakeMerchRowInt{MerchID: 7, Price: 100})

	mockTx := new(Tx)
	mockDB.
		On("Begin", context.Background()).
		Return(mockTx, nil)

	mockTx.
		On("Exec", context.Background(),
			"UPDATE users SET balance = balance - $1 WHERE user_id = $2",
			100, 99).
		Return(pgconn.CommandTag{}, nil)

	mockTx.
		On("Exec", context.Background(),
			"INSERT INTO orders (user_id, merch_id) VALUES ($1, $2)",
			99, 7).
		Return(pgconn.CommandTag{}, nil)

	mockTx.
		On("Commit", context.Background()).
		Return(nil)

	mockTx.
		On("Rollback", context.Background()).
		Return(nil)

	db.DBConn = mockDB

	err := handlers.Buy(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Покупка успешно совершена")
	}
}

func TestBuy_InsufficientFunds(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/buy/cup", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("item")
	c.SetParamValues("cup")
	c.Set("Balance", 50)
	c.Set("userID", 99)

	mockDB := new(mocks.DB)

	mockDB.
		On("QueryRow", context.Background(),
			"SELECT merch_id, price FROM merch WHERE name = $1", "cup").
		Return(&FakeMerchRowInt{MerchID: 7, Price: 100})

	db.DBConn = mockDB

	err := handlers.Buy(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Недостаточно средств")
	}
}
