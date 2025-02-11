package main

import (
	"github.com/prorok210/AvitoShop/internal/db"
	ahs "github.com/prorok210/AvitoShop/internal/services/avitoShopService"
)

func main() {
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	ahs.StartServer()
}
