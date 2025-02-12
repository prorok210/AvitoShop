package main

import (
	"github.com/prorok210/AvitoShop/internal/db"
	ahs "github.com/prorok210/AvitoShop/internal/services/avitoShopService"
)

// @title API Avito shop
// @version 1.0.0
// @description API для отбора на стажировку в Авито
// @host localhost:8083
// @BasePath /api

// @host localhost:8083

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	ahs.StartServer()
}
