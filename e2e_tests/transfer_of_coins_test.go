package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_TransferSuccessfully(t *testing.T) {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	defer db.DBConn.Close()

	e := setupTestServer()
	defer e.Close()

	u1 := &models.User{
		Name:     "John Doe",
		Password: "password123",
	}

	u2 := &models.User{
		Name:     "Jane Doe",
		Password: "125password",
	}

	// Получение токенов
	authReq1 := httptest.NewRequest(http.MethodPost, "/api/auth",
		bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			u1.Name, u1.Password)))
	authReq1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	authRec1 := httptest.NewRecorder()
	authReq2 := httptest.NewRequest(http.MethodPost, "/api/auth",
		bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			u2.Name, u2.Password)))
	authReq2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	authRec2 := httptest.NewRecorder()

	e.ServeHTTP(authRec1, authReq1)
	assert.Equal(t, http.StatusCreated, authRec1.Code)
	var authResp1 models.AuthResponse
	err := json.Unmarshal(authRec1.Body.Bytes(), &authResp1)
	assert.NoError(t, err)
	token1 := authResp1.Token

	e.ServeHTTP(authRec2, authReq2)
	assert.Equal(t, http.StatusCreated, authRec2.Code)
	var authResp2 models.AuthResponse
	err = json.Unmarshal(authRec2.Body.Bytes(), &authResp2)
	assert.NoError(t, err)
	token2 := authResp2.Token

	// Перевод монет от первого пользователя ко второму
	transferReq := httptest.NewRequest(http.MethodPost, "/api/sendCoin",
		bytes.NewBufferString(fmt.Sprintf(`{"toUser": "%s", "amount": 350}`, u2.Name)))
	transferReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	transferReq.Header.Set("Authorization", "Bearer "+token1)
	transferRec1 := httptest.NewRecorder()

	e.ServeHTTP(transferRec1, transferReq)
	assert.Equal(t, http.StatusOK, transferRec1.Code)

	// Проверка баланаса и переводов
	infoReq1 := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq1.Header.Set("Authorization", "Bearer "+token1)
	infoRec1 := httptest.NewRecorder()

	e.ServeHTTP(infoRec1, infoReq1)
	assert.Equal(t, http.StatusOK, infoRec1.Code)
	var infoResp1 models.InfoResponse
	err = json.Unmarshal(infoRec1.Body.Bytes(), &infoResp1)
	assert.NoError(t, err)
	assert.Equal(t, 650, infoResp1.Coins)
	var found bool
	for _, v := range infoResp1.CoinHistory.Sent {
		if v.ToUser == u2.Name {
			found = true
			assert.Equal(t, v.Amount, 350)
		}
	}
	assert.Equal(t, true, found)

	infoReq2 := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq2.Header.Set("Authorization", "Bearer "+token2)
	infoRec2 := httptest.NewRecorder()

	e.ServeHTTP(infoRec2, infoReq2)
	assert.Equal(t, http.StatusOK, infoRec2.Code)
	var infoResp2 models.InfoResponse
	err = json.Unmarshal(infoRec2.Body.Bytes(), &infoResp2)
	assert.NoError(t, err)
	assert.Equal(t, 1350, infoResp2.Coins)
	found = false
	for _, v := range infoResp2.CoinHistory.Received {
		if v.FromUser == u1.Name {
			found = true
			assert.Equal(t, v.Amount, 350)
		}
	}
	assert.Equal(t, true, found)

	// Перевод монет от второго пользователя к первому
	transferReq = httptest.NewRequest(http.MethodPost, "/api/sendCoin",
		bytes.NewBufferString(fmt.Sprintf(`{"toUser": "%s", "amount": 1000}`, u1.Name)))
	transferReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	transferReq.Header.Set("Authorization", "Bearer "+token2)
	transferRec2 := httptest.NewRecorder()

	e.ServeHTTP(transferRec2, transferReq)
	assert.Equal(t, http.StatusOK, transferRec2.Code)

	// Проверка баланаса и переводов
	infoReq1 = httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq1.Header.Set("Authorization", "Bearer "+token1)
	infoRec1 = httptest.NewRecorder()

	e.ServeHTTP(infoRec1, infoReq1)
	assert.Equal(t, http.StatusOK, infoRec1.Code)
	err = json.Unmarshal(infoRec1.Body.Bytes(), &infoResp1)
	assert.NoError(t, err)
	assert.Equal(t, 1650, infoResp1.Coins)
	found = false
	for _, v := range infoResp1.CoinHistory.Received {
		if v.FromUser == u2.Name {
			found = true
			assert.Equal(t, v.Amount, 1000)
		}
	}
	assert.Equal(t, true, found)

	infoReq2 = httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq2.Header.Set("Authorization", "Bearer "+token2)
	infoRec2 = httptest.NewRecorder()

	e.ServeHTTP(infoRec2, infoReq2)
	assert.Equal(t, http.StatusOK, infoRec2.Code)
	err = json.Unmarshal(infoRec2.Body.Bytes(), &infoResp2)
	assert.NoError(t, err)
	assert.Equal(t, 350, infoResp2.Coins)
	found = false
	for _, v := range infoResp2.CoinHistory.Sent {
		if v.ToUser == u1.Name {
			found = true
			assert.Equal(t, v.Amount, 1000)
		}
	}
	assert.Equal(t, true, found)
}

func Test_TransferInsufficientFunds(t *testing.T) {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	defer db.DBConn.Close()

	e := setupTestServer()
	defer e.Close()

	u1 := &models.User{
		Name:     "Walter White",
		Password: "password123",
	}

	u2 := &models.User{
		Name:     "Walter Whitman",
		Password: "228356123",
	}

	// Получение токенов
	authReq1 := httptest.NewRequest(http.MethodPost, "/api/auth",
		bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			u1.Name, u1.Password)))
	authReq1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	authRec1 := httptest.NewRecorder()

	authReq2 := httptest.NewRequest(http.MethodPost, "/api/auth",
		bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			u2.Name, u2.Password)))
	authReq2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	authRec2 := httptest.NewRecorder()

	e.ServeHTTP(authRec1, authReq1)
	assert.Equal(t, http.StatusCreated, authRec1.Code)
	var authResp1 models.AuthResponse
	err := json.Unmarshal(authRec1.Body.Bytes(), &authResp1)
	assert.NoError(t, err)
	token1 := authResp1.Token

	e.ServeHTTP(authRec2, authReq2)
	assert.Equal(t, http.StatusCreated, authRec2.Code)
	var authResp2 models.AuthResponse
	err = json.Unmarshal(authRec2.Body.Bytes(), &authResp2)
	assert.NoError(t, err)
	token2 := authResp2.Token

	// Перевод отрицательного количества монет
	transferReq := httptest.NewRequest(http.MethodPost, "/api/sendCoin",
		bytes.NewBufferString(fmt.Sprintf(`{"toUser": "%s", "amount": -100}`, u2.Name)))
	transferReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	transferReq.Header.Set("Authorization", "Bearer "+token1)
	transferRec1 := httptest.NewRecorder()

	e.ServeHTTP(transferRec1, transferReq)
	assert.Equal(t, http.StatusBadRequest, transferRec1.Code)

	// Перевод монет самому себе
	transferReq = httptest.NewRequest(http.MethodPost, "/api/sendCoin",
		bytes.NewBufferString(fmt.Sprintf(`{"toUser": "%s", "amount": 100}`, u1.Name)))
	transferReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	transferReq.Header.Set("Authorization", "Bearer "+token1)
	transferRec1 = httptest.NewRecorder()

	e.ServeHTTP(transferRec1, transferReq)
	assert.Equal(t, http.StatusBadRequest, transferRec1.Code)

	// Перевод монет без токена
	transferReq = httptest.NewRequest(http.MethodPost, "/api/sendCoin",
		bytes.NewBufferString(fmt.Sprintf(`{"toUser": "%s", "amount": 100}`, u1.Name)))
	transferReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	transferRec1 = httptest.NewRecorder()

	e.ServeHTTP(transferRec1, transferReq)
	assert.Equal(t, http.StatusUnauthorized, transferRec1.Code)

	// Недостачно средств
	transferReq = httptest.NewRequest(http.MethodPost, "/api/sendCoin",
		bytes.NewBufferString(fmt.Sprintf(`{"toUser": "%s", "amount": 5000}`, u1.Name)))
	transferReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	transferReq.Header.Set("Authorization", "Bearer "+token2)
	transferRec1 = httptest.NewRecorder()

	e.ServeHTTP(transferRec1, transferReq)
	assert.Equal(t, http.StatusBadRequest, transferRec1.Code)

	// Проверка инвентарей
	infoReq1 := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq1.Header.Set("Authorization", "Bearer "+token1)
	infoRec1 := httptest.NewRecorder()

	e.ServeHTTP(infoRec1, infoReq1)
	assert.Equal(t, http.StatusOK, infoRec1.Code)
	var infoResp1 models.InfoResponse
	err = json.Unmarshal(infoRec1.Body.Bytes(), &infoResp1)
	assert.NoError(t, err)
	assert.Equal(t, 1000, infoResp1.Coins)
	assert.Equal(t, 0, len(infoResp1.CoinHistory.Sent))
	assert.Equal(t, 0, len(infoResp1.CoinHistory.Received))

	infoReq2 := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq2.Header.Set("Authorization", "Bearer "+token2)
	infoRec2 := httptest.NewRecorder()

	e.ServeHTTP(infoRec2, infoReq2)
	assert.Equal(t, http.StatusOK, infoRec2.Code)
	var infoResp2 models.InfoResponse
	err = json.Unmarshal(infoRec2.Body.Bytes(), &infoResp2)
	assert.NoError(t, err)
	assert.Equal(t, 1000, infoResp2.Coins)
	assert.Equal(t, 0, len(infoResp2.CoinHistory.Sent))
	assert.Equal(t, 0, len(infoResp2.CoinHistory.Received))
}
