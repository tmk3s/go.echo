package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"app/usecase"
)

type JobHandler struct {
	JobUseCase usecase.JobUseCase
}

func NewJobHandler(uc usecase.JobUseCase) *JobHandler {
	return &JobHandler{JobUseCase: uc}
}

func (h *JobHandler) GetStatus(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job id")
	}

	job, err := h.JobUseCase.GetById(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}
	return c.JSON(http.StatusOK, job)
}
