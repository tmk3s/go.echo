package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"app/usecase"
)

type CompanyHandler struct {
	usecase.CompanyUseCase
}

type CompanyUpdateParams struct {
	Name string `json:"name"`
}

func NewCompanyHandler(u usecase.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{u}
}

func (h *CompanyHandler) Show(c echo.Context) error {
	companyId := CurrentCompanyId(c)
	company, err := h.CompanyUseCase.GetCompany(companyId)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) Update(c echo.Context) error {
	var params CompanyUpdateParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	companyId := CurrentCompanyId(c)
	if err := h.CompanyUseCase.UpdateCompany(companyId, params.Name); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}
