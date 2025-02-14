package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/db/mocks"
	"github.com/prorok210/AvitoShop/internal/handlers"
	"github.com/prorok210/AvitoShop/internal/models"
	"github.com/prorok210/AvitoShop/internal/utils"
)

func Test_AuthCreateUser(t *testing.T) {
	e := echo.New()
	body := `{"username":"Jon Snow","password":"123445654"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockDB := new(mocks.DB)

	mockDB.On("QueryRow", context.Background(),
		"SELECT user_id, password FROM users WHERE name = $1", "Jon Snow").
		Return(&FakeUserRow{Err: pgx.ErrNoRows})

	mockDB.On("Exec", context.Background(),
		"INSERT INTO users(name, password) VALUES ($1, $2)",
		"Jon Snow", mock.AnythingOfType("string")).
		Return(pgconn.CommandTag{}, nil)

	mockDB.On("QueryRow", context.Background(),
		"SELECT user_id FROM users WHERE name = $1", "Jon Snow").
		Return(&FakeUserRow{UserID: 42})

	mockDB.On("Exec", context.Background(),
		"INSERT INTO tokens(access_token, user_id) VALUES ($1, $2)",
		mock.AnythingOfType("string"), 42).
		Return(pgconn.CommandTag{}, nil)

	db.DBConn = mockDB

	err := handlers.Auth(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}
}

func Test_AuthGetToken(t *testing.T) {
	e := echo.New()
	body := `{"username":"Jon Snow","password":"123445654"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	hashPass, err := utils.HashPassword("123445654")
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)

	mockDB.On("QueryRow", context.Background(),
		"SELECT user_id, password FROM users WHERE name = $1", "Jon Snow").
		Return(&FakeUserRow{UserID: 42, Password: hashPass}, nil)

	mockDB.On("Exec", context.Background(),
		"INSERT INTO tokens(access_token, user_id) VALUES ($1, $2)",
		mock.AnythingOfType("string"), 42).
		Return(pgconn.CommandTag{}, nil)

	db.DBConn = mockDB

	err = handlers.Auth(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}
}

func Test_AuthInvalidRequest(t *testing.T) {
	e := echo.New()
	body := `{"username":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.Auth(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Неверный запрос")
	}
}

func Test_GetInfo_Success(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", 42)
	c.Set("Balance", 500)

	inventoryRows := &FakeInfoRows{
		data: [][]interface{}{
			{"cup", 2},
		},
	}

	sentRows := &FakeInfoRows{
		data: [][]interface{}{
			{"Alice", 50},
		},
	}

	receivedRows := &FakeInfoRows{
		data: [][]interface{}{
			{"Bob", 30},
		},
	}

	mockDB := new(mocks.DB)

	mockDB.On("Query", context.Background(),
		"SELECT m.name, COUNT(*) FROM orders o JOIN merch m ON o.merch_id = m.merch_id WHERE o.user_id = $1 GROUP BY m.name",
		42).
		Return(inventoryRows, nil)

	mockDB.On("Query", context.Background(),
		"SELECT u.name, t.amount_coins FROM transactions t JOIN users u ON t.recipient_id = u.user_id WHERE t.user_id = $1",
		42).
		Return(sentRows, nil)

	mockDB.On("Query", context.Background(),
		"SELECT u.name, t.amount_coins FROM transactions t JOIN users u ON t.user_id = u.user_id WHERE t.recipient_id = $1",
		42).
		Return(receivedRows, nil)

	db.DBConn = mockDB

	err := handlers.GetInfo(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.InfoResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Coins)
	assert.Len(t, resp.Inventory, 1)
	assert.Equal(t, "cup", resp.Inventory[0].Type)
	assert.Equal(t, 2, resp.Inventory[0].Quantity)
	assert.Len(t, resp.CoinHistory.Sent, 1)
	assert.Equal(t, "Alice", resp.CoinHistory.Sent[0].ToUser)
	assert.Equal(t, 50, resp.CoinHistory.Sent[0].Amount)
	assert.Len(t, resp.CoinHistory.Received, 1)
	assert.Equal(t, "Bob", resp.CoinHistory.Received[0].FromUser)
	assert.Equal(t, 30, resp.CoinHistory.Received[0].Amount)
}

func Test_GetInfoDatabaseError(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", 42)
	c.Set("Balance", 500)

	mockDB := new(mocks.DB)

	mockDB.
		On("Query", context.Background(),
			"SELECT m.name, COUNT(*) FROM orders o JOIN merch m ON o.merch_id = m.merch_id WHERE o.user_id = $1 GROUP BY m.name",
			42).
		Return(nil, errors.New("database connection error"))

	db.DBConn = mockDB

	err := handlers.GetInfo(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var resp models.Error500Response
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Ошибка при запросе инвентаря", resp.Error)
	}
}
