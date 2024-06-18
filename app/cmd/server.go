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

	// https://echo.labstack.com/docs/middleware/cors
	// https://zenn.dev/yuyan/books/c6995204c13a83/viewer/712188
	// これとフロントでsign_inするときに, { withCredentials: true }追加
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins: []string{"http://localhost:3000"},
		// これだとCORSエラーになる。何か足りていない・・
		// AllowHeaders: []string{
		// 	echo.HeaderAccessControlAllowHeaders,
		// 	echo.HeaderContentType,
		// 	echo.HeaderContentLength,
		// 	echo.HeaderAcceptEncoding,
		// 	echo.HeaderXCSRFToken,
		// 	echo.HeaderAuthorization,
		// },
	}))

	e.Logger.Info("hello")
	router.SteupRouter(e)
	e.Logger.Fatal(e.Start(":1323"))
}