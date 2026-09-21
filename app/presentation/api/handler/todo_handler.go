package handler

import (
	"fmt"
	"net/http"
	"strconv"

	// "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"app/usecase"
)

type TodoHandler struct {
	usecase.TodoUseCase
}

type TodoPath struct {
	ID string `param:"id"`
}

type TodoCreateParams struct {
	Title string
}

func NewTodoHandler(u usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{u}
}

func (h *TodoHandler) Index(c echo.Context) error {
	fmt.Printf("%s", "call Index")
	userId := CurrentUserId(c)
	todos, err := h.TodoUseCase.GetTodos(userId)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c echo.Context) error {
	fmt.Printf("%s", "call Create")
	var params TodoCreateParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId := CurrentUserId(c)
	err := h.TodoUseCase.AddTodo(userId, params.Title)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *TodoHandler) Complete(c echo.Context) error {
	fmt.Printf("%s", "call Complete")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid todo id")
	}
	userId := CurrentUserId(c)
	err = h.TodoUseCase.DoneTodo(userId, uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *TodoHandler) Delete(c echo.Context) error {
	fmt.Printf("%s", "call Delete")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid todo id")
	}
	userId := CurrentUserId(c)
	err = h.TodoUseCase.DeleteTodo(userId, uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}
