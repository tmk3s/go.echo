package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"app/usecase"
)

type PrefectureHandler struct {
	usecase.PrefectureUseCase
}

func NewPrefectureHandler(u usecase.PrefectureUseCase) *PrefectureHandler {
	return &PrefectureHandler{u}
}

func (h *PrefectureHandler) Index(c echo.Context) error {
	prefectures, err := h.PrefectureUseCase.GetAll()
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, prefectures)
}
