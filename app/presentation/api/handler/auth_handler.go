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
		SigningKey:     signingKey,
		TokenLookup:   "cookie:session",
	}
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var params SignUpRequest
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if params.Email == "" || params.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email and password are required")
	}

	existing, err := h.AuthUseCase.GetUserByEmail(params.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if existing != nil {
		return echo.NewHTTPError(http.StatusConflict, "email already exists")
	}

	user, err := h.AuthUseCase.CreateUser(params.Email, params.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusCreated, user)
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

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour)
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Secure = os.Getenv("APP_ENV") == "production"
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "ok")
}
