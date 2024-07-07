package handler

import (
	"net/http"
	"fmt"
//	"strconv"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo-jwt/v4"
	"github.com/golang-jwt/jwt/v5"

	"app/db"
)

type TodoHandler struct {}

type TodoPath struct {
	ID string `param:"id"`
}

func (h *TodoHandler) Index(c echo.Context) error{
	fmt.Printf("%s", "call index!!")
	id := userIDFromToken(c)
	user := db.FindUserById(id)
	if user.ID == 0 {
		return echo.ErrNotFound
	}
	todos := db.FindTodos(&db.Todo{ UserId: user.ID })
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c echo.Context) error{
	todo := new(db.Todo)
	if err := c.Bind(todo); err != nil {
			return err
	}
	if todo.Title == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid to or message fields",
		}
	}

	id := userIDFromToken(c)
	user := db.FindUserById(id)
	if user.ID == 0 {
		return echo.ErrNotFound
	}
	todo.UserId = id
  db.CreateTodo(todo)

	todos := db.FindTodos(&db.Todo{ UserId: id })
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Complete(c echo.Context) error{
	todo := db.FindTodo(c.Param("id"))
	if todo.ID == 0 {
		return echo.ErrNotFound
	}	

	id := userIDFromToken(c)
	user := db.FindUserById(id)
	if user.ID == 0 {
		return echo.ErrNotFound
	}

  if err := db.UpdateTodo(&todo); err != nil {
		return err
	}
	todos := db.FindTodos(&db.Todo{ UserId: id })
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Delete(c echo.Context) error{
	todo := db.FindTodo(c.Param("id"))
	
	id := userIDFromToken(c)
	user := db.FindUserById(id)
	if user.ID == 0 {
		return echo.ErrNotFound
	}

	if err := db.DeleteTodo(&todo); err != nil {
		return err
	}

	todos := db.FindTodos(&db.Todo{ UserId: id })
	return c.JSON(http.StatusOK, todos)
}

func userIDFromToken(c echo.Context) uint {
	fmt.Print(c.Get("user"))
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	// id, _ := strconv.Atoi(claims.Id)
	fmt.Printf("%s ID desuto!!!!!!!!!!!!!!!!!!!", claims.Id)
	return claims.Id
}