package server

import (
	"go-app/internal/server/books"
	"go-app/internal/server/sauth"
	"go-app/pkg/auth"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
}

func New(secret string) *Server {
	e := echo.New()
	e.Use(middleware.Logger())

	g := e.Group("/api/v1")

	authMiddleware := echojwt.WithConfig(makeConfig(secret))

	sauth.SetupGroup(g)
	books.SetupGroup(g, authMiddleware)

	return &Server{
		echo: e,
	}
}

func makeConfig(secret string) echojwt.Config {
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
