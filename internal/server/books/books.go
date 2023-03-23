package books

import (
	"fmt"
	"go-app/internal/server"
	"go-app/pkg/books"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type BookRouter struct {
	middleware echo.MiddlewareFunc
}

func NewBookGroup(secret string) *BookRouter {
	authMiddleware := echojwt.WithConfig(server.MakeConfig(secret))
	return &BookRouter{
		authMiddleware,
	}
}

func (b *BookRouter) SetupGroup(parentGroup *echo.Group) {
	ng := parentGroup.Group("/books")
	controller := BookController{bookService{}}

	ng.Use(b.middleware) // add guard

	ng.GET("", controller.listBooks)
}

type BookController struct {
	bookSvc books.BookService
}

func (b *BookController) listBooks(c echo.Context) error {

	user := c.Get("user")
	if user != nil {
		x := user.(*jwt.Token)
		fmt.Println(x)
		fmt.Println(c.Get("user"), user)
	}

	// claims := user.Claims.(*auth.JwtCustomClaims)

	return c.JSON(http.StatusOK, b.bookSvc.ListBooks())
}

type bookService struct {
}

func (b bookService) ListBooks() []books.Book {
	return []books.Book{{
		Id:    "book-001",
		Title: "Who steal my book?",
	}}
}
