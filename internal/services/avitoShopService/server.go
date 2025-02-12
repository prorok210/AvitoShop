package avitoShopService

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prorok210/AvitoShop/config"
	h "github.com/prorok210/AvitoShop/internal/handlers"
	"github.com/prorok210/AvitoShop/internal/middlewares"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func StartServer() {
	e := echo.New()
	e.Server.ReadTimeout = config.READTIMEOUT
	e.Server.WriteTimeout = config.WRITETIMEOUT

	e.POST("/api/auth", h.Auth)

	e.GET("/api/buy/:item", h.Buy, middlewares.AuthMiddleware())
	e.POST("api/sendCoin", h.SendCoin, middlewares.AuthMiddleware())
	e.GET("api/info", h.GetInfo, middlewares.AuthMiddleware())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
			echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.Static("/docs", "docs")
	wrapHandler := echoSwagger.EchoWrapHandler(echoSwagger.URL("http://localhost:8083/docs/swagger.json"))
	e.GET("/swagger/*", wrapHandler)

	if err := e.Start(os.Getenv("SERVER_PORT")); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
