package handler

import (
	"errors"
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

// CompanyRegisterParams は未ログインで叩ける会社登録のリクエストボディ。
type CompanyRegisterParams struct {
	CompanyName string `json:"company_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
}

func NewCompanyHandler(u usecase.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{u}
}

// Register は会社と最初のユーザーを作成する。認証前に実行される。
func (h *CompanyHandler) Register(c echo.Context) error {
	var params CompanyRegisterParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	company, err := h.CompanyUseCase.Register(usecase.CompanyRegistrationInput{
		CompanyName: params.CompanyName,
		Email:       params.Email,
		Password:    params.Password,
		LastName:    params.LastName,
		FirstName:   params.FirstName,
	})
	if err != nil {
		// 入力不備だけを利用者に返す。DB エラーの本文は外に出さない。
		var invalid *usecase.ValidationError
		if errors.As(err, &invalid) {
			return echo.NewHTTPError(http.StatusBadRequest, invalid.Message)
		}
		c.Logger().Errorf("company registration failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "登録に失敗しました")
	}
	return c.JSON(http.StatusCreated, company)
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
