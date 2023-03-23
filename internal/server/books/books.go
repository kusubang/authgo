package books

import (
	"fmt"
	"go-app/pkg/auth"
	"go-app/pkg/usr"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type BookRouter struct {
}

func SetupGroup(g *echo.Group, middleware echo.MiddlewareFunc) {
	ng := g.Group("/books")
	authSvc := auth.New(&usr.UserServiceMem{})
	controller := BookController{authSvc}

	// ng.Use(echojwt.WithConfig(config))
	ng.Use(middleware)

	ng.GET("", controller.listBooks)
}

type BookController struct {
	auth *auth.Auth
}

func (b *BookController) listBooks(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*auth.JwtCustomClaims)

	fmt.Println(c.Get("user"), user)

	return c.JSON(http.StatusOK, echo.Map{
		"books": "book1 book2",
	})
}

type BookService struct {
}

func (b *BookService) listBooks() {

}
