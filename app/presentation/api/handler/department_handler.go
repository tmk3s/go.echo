package handler

import (
	"fmt"
	"net/http"

	// "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"app/usecase"
)

type DepartmentHandler struct {
	usecase.DepartmentUseCase
}

type DepartmentPath struct {
	ID string `param:"id"`
}

type DepartmentCreateParams struct {
	Name     string
	ParentId *uint `json:"parent_id"` // null許容のため
}

func NewDepartmentHandler(u usecase.DepartmentUseCase) *DepartmentHandler {
	return &DepartmentHandler{u}
}

func (h *DepartmentHandler) Index(c echo.Context) error {
	fmt.Printf("%s", "call Index")
	companyId := CurrentCompanyId(c)
	departments, err := h.DepartmentUseCase.GetDepartments(companyId)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, departments)
}

func (h *DepartmentHandler) Create(c echo.Context) error {
	fmt.Println("call Create")
	var params DepartmentCreateParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	companyId := CurrentCompanyId(c)
	err := h.DepartmentUseCase.Create(companyId, params.Name, params.ParentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}
