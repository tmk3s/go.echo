package main

import (
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"app/config"
	"app/infrastructure/worker"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000"},
	}))
	// ChromeのPrivate Network Access対応
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Access-Control-Request-Private-Network") == "true" {
				c.Response().Header().Set("Access-Control-Allow-Private-Network", "true")
			}
			return next(c)
		}
	})

	dbConnection, err := config.NewMysqlConnection()
	if err != nil {
		panic(fmt.Sprintf("create dbConnection failed to connect database %s", err))
	}

	err = config.ExecuteMigrate(dbConnection)
	if err != nil {
		panic(fmt.Sprintf("create dbConnection failed to connect database %s", err))
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer asynqClient.Close()

	r := registry.NewReigistry(dbConnection, asynqClient)
	h := r.NewAppHandler()

	// Asynq ワーカーをバックグラウンドで起動
	asynqServer := worker.NewServer(redisAddr)
	mux := worker.NewServeMux(r.NewEmployeeUseCase())
	go func() {
		if err := asynqServer.Run(mux); err != nil {
			e.Logger.Errorf("asynq worker error: %v", err)
		}
	}()

	e.Logger.Info("hello")
	router.SteupRouter(e, *h)
	e.Logger.Fatal(e.Start(":1323"))
}
