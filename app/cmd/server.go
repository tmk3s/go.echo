package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"app/presentation/api/router"
)

func main() {
	e := echo.New()

	// dockerでログ確認するため
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: os.Stdout,
	})
	e.Use(logger)

	e.Logger.Info("hello")
	router.SteupRouter(e)
	e.Logger.Fatal(e.Start(":1323"))
}