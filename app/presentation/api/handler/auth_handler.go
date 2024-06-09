package handler

import (
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-jwt/v4"
    "github.com/golang-jwt/jwt/v5"

	"app/db"
)

type AuthHandler struct {}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
    Id    int    `json:"id"`
    Email string `json:"email"`
    jwt.RegisteredClaims
}

var signingKey = []byte("secret")

var Config = echojwt.Config{
    NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(jwtCustomClaims) },
    SigningKey: signingKey,
}

func (h *AuthHandler) SignUp(c echo.Context) error {
    user := new(db.User)
    if err := c.Bind(user); err != nil {
        return err
    }

    if user.Email == "" || user.Password == "" {
        return &echo.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "invalid Email or password",
        }
    }

    if u := db.FindUser(&db.User{Email: user.Email}); u.Id != 0 {
        return &echo.HTTPError{
            Code:    http.StatusConflict,
            Message: "email already exists",
        }
    }

    db.CreateUser(user)
    user.Password = ""

    return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) SignIn(c echo.Context) error {
    u := new(db.User)
    if err := c.Bind(u); err != nil {
        return err
    }

    user := db.FindUser(&db.User{Email: u.Email})
    if user.Id == 0 || user.Password != u.Password {
        return &echo.HTTPError{
            Code:    http.StatusUnauthorized,
            Message: "invalid Email or password",
        }
    }

    claims := &jwtCustomClaims{
        user.Id,
        user.Email,
        jwt.RegisteredClaims{
            // https://github.com/golang-jwt/jwt/blob/main/example_test.go
            ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
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

    // return c.JSON(http.StatusOK, map[string]string{
    //     "token": t,
    // })
}