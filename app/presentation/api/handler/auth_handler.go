package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"app/usecase"
)

type AuthHandler struct {
	usecase.AuthUseCase
}

func NewAuthHandler(u usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{u}
}

type jwtCustomClaims struct {
	Id        uint   `json:"id"`
	Email     string `json:"email"`
	CompanyId uint   `json:"company_id"`
	jwt.RegisteredClaims
}

var signingKey []byte
var Config echojwt.Config

func init() {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		panic("JWT_SECRET environment variable must be set")
	}
	signingKey = []byte(key)
	Config = echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(jwtCustomClaims) },
		SigningKey:    signingKey,
		TokenLookup:   "cookie:session",
	}
}

// ユーザー単体のサインアップは廃止した。
// 会社に所属しないユーザー(company_id = 0)が作られると、
// 他社のデータが見えたり何も操作できなかったりするため、
// 会社登録(POST /companies)で会社とユーザーを同時に作成する。

const sessionCookieName = "session"

// newSessionCookie はセッションCookieを組み立てる。
// 削除時は発行時と同じ属性(特に Path)でないとブラウザが消してくれないため、
// 発行とログアウトで必ずこの関数を通す。
func newSessionCookie(value string, expires time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   os.Getenv("APP_ENV") == "production",
	}
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	var params SignInRequest
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.AuthUseCase.GetUser(params.Email, params.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	}

	claims := &jwtCustomClaims{
		user.ID,
		user.Email,
		user.CompanyId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	expires := time.Now().Add(time.Hour)
	c.SetCookie(newSessionCookie(t, expires, int(time.Hour.Seconds())))

	return c.String(http.StatusOK, "ok")
}

// SignOut はセッションCookieを失効させる。
// Cookie は HttpOnly なのでフロントエンドからは削除できず、サーバーで消す必要がある。
// また、トークンが期限切れでもログアウトできるよう、JWT 必須のグループには置かない。
func (h *AuthHandler) SignOut(c echo.Context) error {
	// MaxAge を負値にすると、ブラウザは即座に Cookie を削除する
	c.SetCookie(newSessionCookie("", time.Unix(0, 0), -1))
	return c.String(http.StatusOK, "ok")
}
