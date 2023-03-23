package sauth

import (
	"go-app/pkg/auth"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	auth *auth.Auth
}

func (a *AuthController) Login(c echo.Context) error {
	id := c.FormValue("id")
	email := c.FormValue("email")
	pw := c.FormValue("pw")

	at, rt, err := a.auth.Login(id, email, pw)

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

func (a *AuthController) RefreshToken(c echo.Context) error {

	refreshToken, err := c.Cookie("refresh-token")
	if err != nil {
		return err
	}
	at, rt, err := a.auth.RefreshToken(refreshToken.Value)

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
