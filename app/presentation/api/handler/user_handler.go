package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"app/usecase"
)

// https://zenn.dev/keitakn/articles/go-naming-rules
// 関数、type、構造体
// キャメルケースで命名します。
// 外部に公開する関数や構造体の場合、先頭を大文字にするという言語仕様があるので、それに合わせてアッパーキャメルケース（先頭大文字から始まる）かローワーキャメルケース（先頭小文字から始まる）が決まります。

type UserHandler struct {
	usecase.UserUseCase
}

func NewUserHandler(u usecase.UserUseCase) *UserHandler {
	return &UserHandler{u}
}

func (h *UserHandler) Index(c echo.Context) error {
	fmt.Printf("%s", "call Index")
	userId := CurrentUserId(c)
	user, err := h.UserUseCase.Get(userId)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c echo.Context) error {
	fmt.Printf("%s", "call Update")
	userId := CurrentUserId(c)
	var params usecase.UserParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.UserUseCase.Update(userId, params)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, nil)
}
