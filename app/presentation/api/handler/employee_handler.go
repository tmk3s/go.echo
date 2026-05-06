package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"app/usecase"
)

type EmployeeHandler struct {
	usecase.EmployeeUseCase
}

type EmployeeCreateParams struct {
	LastName      string `json:"last_name"`
	FirstName     string `json:"first_name"`
	LastNameKana  string `json:"last_name_kana"`
	FirstNameKana string `json:"first_name_kana"`
	Email         string `json:"email"`
	StaffCode     string `json:"staff_code"`
}

type TenureUpdateParams struct {
	ID              uint    `json:"id"`
	JoinedOn        string  `json:"joined_on"`
	ResignationOn   *string `json:"resignation_on"`
	ResignationType *string `json:"resignation_type"`
	Status          *string `json:"status"`
}

type EmployeeUpdateAllParams struct {
	LastName      string               `json:"last_name"`
	FirstName     string               `json:"first_name"`
	LastNameKana  string               `json:"last_name_kana"`
	FirstNameKana string               `json:"first_name_kana"`
	Email         string               `json:"email"`
	StaffCode     string               `json:"staff_code"`
	PostCode      *string              `json:"post_code"`
	PrefectureId  *uint                `json:"prefecture_id"`
	City          *string              `json:"city"`
	AddressLine1  *string              `json:"address_line1"`
	AddressLine2  *string              `json:"address_line2"`
	Tel           *string              `json:"tel"`
	Tenures       []TenureUpdateParams `json:"tenures"`
	DepartmentIds []uint               `json:"department_ids"`
}

func NewEmployeeHandler(u usecase.EmployeeUseCase) *EmployeeHandler {
	return &EmployeeHandler{u}
}

func (h *EmployeeHandler) Index(c echo.Context) error {
	companyId := CurrentCompanyId(c)
	employees, err := h.EmployeeUseCase.GetEmployees(companyId)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) Show(c echo.Context) error {
	companyId := CurrentCompanyId(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	detail, err := h.EmployeeUseCase.GetEmployeeDetail(companyId, uint(id))
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *EmployeeHandler) Create(c echo.Context) error {
	var params EmployeeCreateParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	companyId := CurrentCompanyId(c)
	err := h.EmployeeUseCase.Create(companyId, params.LastName, params.FirstName, params.LastNameKana, params.FirstNameKana, params.Email, params.StaffCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *EmployeeHandler) Export(c echo.Context) error {
	companyId := CurrentCompanyId(c)
	data, err := h.EmployeeUseCase.ExportCSV(companyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\"employees.csv\"")
	return c.Blob(http.StatusOK, "text/csv; charset=utf-8", data)
}

func (h *EmployeeHandler) BulkCreate(c echo.Context) error {
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer file.Close()

	companyId := CurrentCompanyId(c)
	if err := h.EmployeeUseCase.BulkCreateFromCSV(companyId, file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *EmployeeHandler) BulkUpdate(c echo.Context) error {
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer file.Close()

	companyId := CurrentCompanyId(c)
	if err := h.EmployeeUseCase.BulkUpdateFromCSV(companyId, file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *EmployeeHandler) Update(c echo.Context) error {
	var params EmployeeUpdateAllParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	companyId := CurrentCompanyId(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tenures := make([]usecase.TenureUpdateInput, len(params.Tenures))
	for i, t := range params.Tenures {
		tenures[i] = usecase.TenureUpdateInput{
			ID:              t.ID,
			JoinedOn:        t.JoinedOn,
			ResignationOn:   t.ResignationOn,
			ResignationType: t.ResignationType,
			Status:          t.Status,
		}
	}

	input := usecase.UpdateAllInput{
		LastName:      params.LastName,
		FirstName:     params.FirstName,
		LastNameKana:  params.LastNameKana,
		FirstNameKana: params.FirstNameKana,
		Email:         params.Email,
		StaffCode:     params.StaffCode,
		PostCode:      params.PostCode,
		PrefectureId:  params.PrefectureId,
		City:          params.City,
		AddressLine1:  params.AddressLine1,
		AddressLine2:  params.AddressLine2,
		Tel:           params.Tel,
		Tenures:       tenures,
		DepartmentIds: params.DepartmentIds,
	}

	if err := h.EmployeeUseCase.UpdateAll(companyId, uint(id), input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}
