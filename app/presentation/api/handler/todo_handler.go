package handler

import (
	"net/http"
	"fmt"

	"github.com/labstack/echo/v4"
	// "github.com/golang-jwt/jwt/v5"

	// "app/db"
	"app/usecase"
)

type TodoHandler struct {
	usecase.TodoUseCase
}

type TodoPath struct {
	ID string `param:"id"`
}

func (h *TodoHandler) Index(c echo.Context) error{
	fmt.Printf("%s", "sumihisa tomoki")
	userId := CurrentUserId(c)
	todos, err := h.TodoUseCase.GetTodos(userId)
	if err != nil {
			return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c echo.Context) error{
	// todo := new(db.Todo)
	// if err := c.Bind(todo); err != nil {
	// 		return err
	// }
	// if todo.Title == "" {
	// 	return &echo.HTTPError{
	// 		Code:    http.StatusBadRequest,
	// 		Message: "invalid to or message fields",
	// 	}
	// }

	// id := userIDFromToken(c)
	// user := db.FindUserById(id)
	// if user.ID == 0 {
	// 	return echo.ErrNotFound
	// }
	// todo.UserId = id
  // db.CreateTodo(todo)

	// todos := db.FindTodos(&db.Todo{ UserId: id })
	// return c.JSON(http.StatusOK, todos)
	return c.JSON(http.StatusOK, nil)
}

func (h *TodoHandler) Complete(c echo.Context) error{
	// todo := db.FindTodo(c.Param("id"))s
	return c.JSON(http.StatusOK, nil)
}

func (h *TodoHandler) Delete(c echo.Context) error{
	// todo := db.FindTodo(c.Param("id"))
	
	// id := userIDFromToken(c)
	// user := db.FindUserById(id)
	// if user.ID == 0 {
	// 	return echo.ErrNotFound
	// }

	// if err := db.DeleteTodo(&todo); err != nil {
	// 	return err
	// }

	// todos := db.FindTodos(&db.Todo{ UserId: id })
	// return c.JSON(http.StatusOK, todos)
	return c.JSON(http.StatusOK, nil)
}