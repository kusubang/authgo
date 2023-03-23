package sauth

import (
	"go-app/pkg/auth"
	"go-app/pkg/usr"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
}

func SetupGroup(g *echo.Group) {
	ng := g.Group("/auth")
	authSvc := auth.New(&usr.UserServiceMem{})
	authCtl := AuthController{authSvc}

	ng.POST("/login", authCtl.Login)
	ng.GET("/refresh", authCtl.RefreshToken)
}
