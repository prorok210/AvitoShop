package e2e_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/internal/db"
	"github.com/prorok210/AvitoShop/internal/handlers"
	"github.com/prorok210/AvitoShop/internal/middlewares"
	"github.com/prorok210/AvitoShop/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *echo.Echo {
	e := echo.New()
	e.POST("/api/auth", handlers.Auth)
	e.GET("/api/buy/:item", handlers.Buy, middlewares.AuthMiddleware())
	e.GET("/api/info", handlers.GetInfo, middlewares.AuthMiddleware())
	e.POST("/api/sendCoin", handlers.SendCoin, middlewares.AuthMiddleware())
	return e
}

func Test_FullBuyScenario(t *testing.T) {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	defer db.DBConn.Close()

	e := setupTestServer()
	defer e.Close()

	// Креды для регистрации
	authJSON := `{
        "username": "Jeff Smith",
        "password": "password123"
    }`

	// Получние токена
	authReq := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(authJSON))
	authReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	authRec := httptest.NewRecorder()

	e.ServeHTTP(authRec, authReq)
	assert.Equal(t, http.StatusCreated, authRec.Code)

	var authResp models.AuthResponse
	err := json.Unmarshal(authRec.Body.Bytes(), &authResp)
	assert.NoError(t, err)
	token := authResp.Token

	// Получение информации о пользователе
	infoReq := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq.Header.Set("Authorization", "Bearer "+token)
	infoRec := httptest.NewRecorder()

	e.ServeHTTP(infoRec, infoReq)
	assert.Equal(t, http.StatusOK, infoRec.Code)

	var initialInfo models.InfoResponse
	err = json.Unmarshal(infoRec.Body.Bytes(), &initialInfo)
	assert.NoError(t, err)
	initialBalance := initialInfo.Coins

	// Покупка кружки
	buyReq := httptest.NewRequest(http.MethodGet, "/api/buy/cup", nil)
	buyReq.Header.Set("Authorization", "Bearer "+token)
	buyRec := httptest.NewRecorder()

	e.ServeHTTP(buyRec, buyReq)
	assert.Equal(t, http.StatusOK, buyRec.Code)

	finalInfoReq := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	finalInfoReq.Header.Set("Authorization", "Bearer "+token)
	finalInfoRec := httptest.NewRecorder()

	e.ServeHTTP(finalInfoRec, finalInfoReq)
	assert.Equal(t, http.StatusOK, finalInfoRec.Code)

	var finalInfo models.InfoResponse
	err = json.Unmarshal(finalInfoRec.Body.Bytes(), &finalInfo)
	assert.NoError(t, err)
	//Проверка, что баланс уменьшился
	assert.Less(t, finalInfo.Coins, initialBalance)

	// Проверка, что кружка появилась в инвентаре
	found := false
	for _, item := range finalInfo.Inventory {
		if item.Type == "cup" {
			found = true
			assert.GreaterOrEqual(t, item.Quantity, 1)
			break
		}
	}
	assert.True(t, found, "Купленный товар не найден в инвентаре")
}

func TestBuyScenario_InsufficientFunds(t *testing.T) {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	defer db.DBConn.Close()

	e := setupTestServer()
	defer e.Close()

	// Креды для регистрации
	authJSON := `{
        "username": "Poor John",
        "password": "password123"
    }`
	// Получние токена
	authReq := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(authJSON))
	authReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	authRec := httptest.NewRecorder()

	e.ServeHTTP(authRec, authReq)
	assert.Equal(t, http.StatusCreated, authRec.Code)

	var authResp models.AuthResponse
	err := json.Unmarshal(authRec.Body.Bytes(), &authResp)
	assert.NoError(t, err)
	token := authResp.Token

	// Получение информации о пользователе
	infoReq := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	infoReq.Header.Set("Authorization", "Bearer "+token)
	infoRec := httptest.NewRecorder()

	e.ServeHTTP(infoRec, infoReq)
	assert.Equal(t, http.StatusOK, infoRec.Code)

	var initialInfo models.InfoResponse
	err = json.Unmarshal(infoRec.Body.Bytes(), &initialInfo)
	assert.NoError(t, err)
	assert.Equal(t, 1000, initialInfo.Coins)

	// Пытаемся купить 3 худи
	for i := 0; i < 2; i++ {
		buyReq := httptest.NewRequest(http.MethodGet, "/api/buy/pink-hoody", nil)
		buyReq.Header.Set("Authorization", "Bearer "+token)
		buyRec := httptest.NewRecorder()

		// Первые 2 раза покупка проходит успешно
		e.ServeHTTP(buyRec, buyReq)
		assert.Equal(t, http.StatusOK, buyRec.Code)
	}

	// 3 раз - недостаточно средств
	buyReq := httptest.NewRequest(http.MethodGet, "/api/buy/pink-hoody", nil)
	buyReq.Header.Set("Authorization", "Bearer "+token)
	buyRec := httptest.NewRecorder()

	e.ServeHTTP(buyRec, buyReq)
	assert.Equal(t, http.StatusBadRequest, buyRec.Code)

	var errorResp models.Error400Response
	err = json.Unmarshal(buyRec.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Contains(t, errorResp.Error, "Недостаточно средств")

	finalInfoReq := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	finalInfoReq.Header.Set("Authorization", "Bearer "+token)
	finalInfoRec := httptest.NewRecorder()

	e.ServeHTTP(finalInfoRec, finalInfoReq)
	assert.Equal(t, http.StatusOK, finalInfoRec.Code)

	var finalInfo models.InfoResponse
	err = json.Unmarshal(finalInfoRec.Body.Bytes(), &finalInfo)
	assert.NoError(t, err)
	assert.Equal(t, 0, finalInfo.Coins)

	// Проверяем количество худи в инвентаре
	found := false
	for _, item := range finalInfo.Inventory {
		if item.Type == "pink-hoody" {
			found = true
			assert.Equal(t, 2, item.Quantity, "В инвентаре должно быть ровно 2 худи")
			break
		}
	}
	assert.True(t, found, "Худи не найдены в инвентаре")
}
