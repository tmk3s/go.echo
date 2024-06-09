package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// https://zenn.dev/keitakn/articles/go-naming-rules
// 関数、type、構造体
// キャメルケースで命名します。
// 外部に公開する関数や構造体の場合、先頭を大文字にするという言語仕様があるので、それに合わせてアッパーキャメルケース（先頭大文字から始まる）かローワーキャメルケース（先頭小文字から始まる）が決まります。

type UserHandler struct {}

func (h *UserHandler) Index(c echo.Context) error{
	return c.String(http.StatusOK, "User INdex!OKOKOKO")
}