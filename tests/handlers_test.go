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

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/db/mocks"
	"github.com/prorok210/AvitoShop/internal/handlers"
)

func TestAuth_CreateUser(t *testing.T) {
	e := echo.New()
	body := `{"name":"Jon Snow","email":"jon@snow.com","password":"123445654"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockDB := new(mocks.DB)

	mockDB.On("QueryRow", context.Background(),
		"SELECT user_id, password FROM users WHERE email = $1", "jon@snow.com").
		Return(&fakeRow{Err: pgx.ErrNoRows})

	mockDB.On("Exec", context.Background(),
		"INSERT INTO users(name, email, password) VALUES ($1, $2, $3)",
		"Jon Snow", "jon@snow.com", mock.AnythingOfType("string")).
		Return(pgconn.CommandTag{}, nil)

	mockDB.On("QueryRow", context.Background(),
		"SELECT user_id FROM users WHERE email = $1", "jon@snow.com").
		Return(&fakeRow{UserID: 42})

	mockDB.On("Exec", context.Background(),
		"INSERT INTO tokens(access_token, refresh_token, user_id) VALUES ($1, $2, $3)",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), 42).
		Return(pgconn.CommandTag{}, nil)

	db.DBConn = mockDB

	err := handlers.Auth(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "access_token")
	}
}

type fakeRow struct {
	UserID int
	Pwd    string
	Err    error
}

func (r *fakeRow) Scan(dest ...interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	if len(dest) >= 1 {
		if ref, ok := dest[0].(*int); ok {
			*ref = r.UserID
		}
	}
	if len(dest) >= 2 {
		if ref, ok := dest[1].(*string); ok {
			*ref = r.Pwd
		}
	}
	return nil
}
