package main

import (
	"fmt"
	"go-app/pkg/auth"
	"go-app/pkg/usr"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Sub   string `json:"sub"`
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func MakeLogin(auth *auth.Auth) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.FormValue("id")
		email := c.FormValue("email")
		pw := c.FormValue("pw")
		_, err := auth.UserService.Find(id, pw)

		if err != nil {
			return echo.ErrUnauthorized
		}

		at, rt, err := auth.Login(id, email, pw)

		if err != nil {
			return echo.ErrBadRequest
		}

		//
		cookie := new(http.Cookie)
		cookie.Name = "refresh-token"
		cookie.Value = rt
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, echo.Map{
			"accessToken": at,
			// "refreshToken": rt,
		})

	}
}

func readCookie(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return c.String(http.StatusOK, "read a cookie")
}

func accessible(c echo.Context) error {
	readCookie(c)
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Id
	fmt.Println(claims)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

const SECRET = "secret"

func RefreshToken(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh-token")
	if err != nil {
		return err
	}
	// fmt.Println("refresh token:", refreshToken.String(), refreshToken.Value)

	token, err := jwt.Parse(refreshToken.Value, func(token *jwt.Token) (interface{}, error) {
		// 클라이언트로 받은 토큰이 HMAC 알고리즘이 맞는지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(claims, ok)
	return nil
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	auth := auth.New(&usr.UserServiceMem{})
	loginHander := MakeLogin(auth)

	e.GET("/", func(c echo.Context) error {
		readCookie(c)
		return c.String(http.StatusOK, "HELLO WORLD!")
	}) // Login route

	e.POST("/login", loginHander)

	// Unauthenticated route
	e.GET("/", accessible)
	e.GET("/refresh", RefreshToken)
	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(SECRET),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":1323"))

}
