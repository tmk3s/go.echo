package handler

import (
    "net/http"
    "fmt"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-jwt/v4"
    "github.com/golang-jwt/jwt/v5"

    "app/usecase"
)

type AuthHandler struct {
    usecase.AuthUseCase
}

func NewAuthHandler(u usecase.AuthUseCase) *AuthHandler {
    return &AuthHandler{u}
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
    Id    uint    `json:"id"`
    Email string  `json:"email"`
    jwt.RegisteredClaims
}

var signingKey = []byte("secret")

var Config = echojwt.Config{
    NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(jwtCustomClaims) },
    SigningKey: signingKey,
}

type SignUpRequest struct {
    Email string `json:email`
    Password string `json:password`
}

func (h *AuthHandler) SignUp(c echo.Context) error {
    var params SignUpRequest
    if err := c.Bind(&params); err != nil {
        return err
    }

    fmt.Printf("%s", "Bind Done")

    if params.Email == "" || params.Password == "" {
        return &echo.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "invalid Email or password",
        }
    }

    user, err := h.AuthUseCase.GetUser(params.Email, params.Password)
    if err != nil {
        return err
    }

    if user != nil {
        return &echo.HTTPError{
            Code:    http.StatusConflict,
            Message: "email already exists",
        }
    }

    user, err = h.AuthUseCase.CreateUser(params.Email, params.Password)
    if err != nil {
        return err
    }

    user.Password = "" // フロントに返さないようにクリア
    return c.JSON(http.StatusCreated, user)
}

type SignInRequest struct {
    Email string `json:email`
    Password string `json:password`
}

func (h *AuthHandler) SignIn(c echo.Context) error {
    fmt.Printf("%s", "SignIn Start")
    var params SignInRequest
    if err := c.Bind(&params); err != nil {
        return err
    }
    fmt.Printf("%s", "Bind Done")
    fmt.Printf("%s", params.Email)
    fmt.Printf("%s", params.Password)

    user, err := h.AuthUseCase.GetUser(params.Email, params.Password)
    if err != nil {
        return err
    }

    if user != nil {
        return &echo.HTTPError{
            Code:    http.StatusUnauthorized,
            Message: "invalid Email or password",
        }
    }

    claims := &jwtCustomClaims{
        user.ID,
        user.Email,
        jwt.RegisteredClaims{
            // https://github.com/golang-jwt/jwt/blob/main/example_test.go
            ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Add(time.Hour * 72).Unix(), 0)),
            // ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72).Unix()), NG
            // ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), NG
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString(signingKey)
    if err != nil {
        return err
    }

    cookie := new(http.Cookie)
    cookie.Name = "session"
    cookie.Value = t
    cookie.Expires = time.Now().Add(24 * time.Hour)
    c.SetCookie(cookie)

    return c.String(http.StatusOK, "write a cookie")
}