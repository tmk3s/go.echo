package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"app/config"
	"app/presentation/api/router"
	"app/registry"
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
		AllowOrigins:     []string{"http://localhost:3000"},
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

	dbConnection, err := config.NewMysqlConnection()
	if err != nil {
		panic(fmt.Sprintf("create dbConnection failed to connect database %s", err))
	}

	err = config.ExecuteMigrate(dbConnection)
	if err != nil {
		panic(fmt.Sprintf("create dbConnection failed to connect database %s", err))
	}

	r := registry.NewReigistry(dbConnection)
	h := r.NewAppHandler()

	e.Logger.Info("hello")
	router.SteupRouter(e, *h)
	e.Logger.Fatal(e.Start(":1323"))
}
