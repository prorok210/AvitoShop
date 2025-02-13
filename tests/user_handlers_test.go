package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// Echo + Aссерт
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/db/mocks"
	"github.com/prorok210/AvitoShop/internal/handlers"
	"github.com/prorok210/AvitoShop/internal/utils"
)

func TestAuth_CreateUser(t *testing.T) {
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

func TestAuth_GetToken(t *testing.T) {
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

	mockDB.
		On("QueryRow", context.Background(),
			"SELECT user_id, password FROM users WHERE name = $1", "Jon Snow").
		Return(&FakeUserRow{UserID: 42, Password: hashPass}, nil)

	mockDB.
		On("Exec", context.Background(),
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

func TestAuth_InvalidRequest(t *testing.T) {
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
