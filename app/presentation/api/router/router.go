package router

// https://qiita.com/ogady/items/0cedd3599c4dc13e9a95 絶対パスの方がいいらしい
import (
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-jwt/v4"
	// "github.com/golang-jwt/jwt/v5"

	"app/presentation/api/handler"
)

func SteupRouter(e *echo.Echo, h handler.AppHandler) {
	// https://echo.labstack.com/docs/routing
	// e.GET("/users", new(handler.UserHandler).Index)

	//　todoのような書き方できるなら無理にclass使わなくても・・・
	e.POST("/sign_up", h.AuthHandler.SignUp)
	e.POST("/sign_in", h.AuthHandler.SignIn)

	api := e.Group("/api")
	api.Use(echojwt.WithConfig(handler.Config)) // /api 下はJWTの認証が必要
	api.GET("/todos", h.TodoHandler.Index)
	api.POST("/todo", h.TodoHandler.Create)
	api.DELETE("/todo/:id", h.TodoHandler.Delete)
	api.PUT("/todo/:id/completed", h.TodoHandler.Complete)

	api.GET("/departments", h.DepartmentHandler.Index)
	api.POST("/department", h.DepartmentHandler.Create)

	api.GET("/user", h.UserHandler.Index)
	api.PUT("/user", h.UserHandler.Update)
}
