package auth

import (
	"go-app/pkg/auth"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	authSvc *auth.Auth
}

func NewAuthGroup(svc *auth.Auth) *AuthRouter {
	return &AuthRouter{
		svc,
	}
}

func (a *AuthRouter) SetupGroup(g *echo.Group) {
	ng := g.Group("/auth")
	authCtl := AuthController{a.authSvc}

	ng.POST("/login", authCtl.login)
	ng.GET("/refresh", authCtl.refreshToken)
}

type AuthController struct {
	authSvc *auth.Auth
}

func (a *AuthController) login(c echo.Context) error {
	id := c.FormValue("id")
	email := c.FormValue("email")
	pw := c.FormValue("pw")

	at, rt, err := a.authSvc.Login(id, email, pw)

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

func (a *AuthController) refreshToken(c echo.Context) error {

	refreshToken, err := c.Cookie("refresh-token")
	if err != nil {
		return err
	}
	at, rt, err := a.authSvc.RefreshToken(refreshToken.Value)

	if err != nil {
		return echo.ErrBadRequest
	}

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
