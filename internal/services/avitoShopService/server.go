package avitoShopService

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/prorok210/AvitoShop/config"
	h "github.com/prorok210/AvitoShop/internal/handlers"
)

func StartServer() {
	e := echo.New()
	e.Server.ReadTimeout = config.READTIMEOUT
	e.Server.WriteTimeout = config.WRITETIMEOUT

	e.POST("/api/auth", h.Auth)

	if err := e.Start(os.Getenv("SERVER_PORT")); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
