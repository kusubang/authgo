package server

import (
	"go-app/pkg/auth"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	g    *echo.Group
}

func New(apiRoot *string, secret string) *Server {
	e := echo.New()
	e.Use(middleware.Logger())

	var prefix string = "/api/v1"
	if apiRoot != nil {
		prefix = *apiRoot
	}

	gAPI := e.Group(prefix)

	return &Server{
		echo: e,
		g:    gAPI,
	}
}

type RouterGroup interface {
	SetupGroup(g *echo.Group)
}

func (s *Server) AddGroup(rg RouterGroup) {
	rg.SetupGroup(s.g)
}

func MakeConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
}

func (s *Server) Start(addr string) error {
	return s.echo.Start(addr)
}
