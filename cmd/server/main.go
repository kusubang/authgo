package main

import (
	"errors"
	"go-app/internal/server"
	iauth "go-app/internal/server/auth"
	"go-app/internal/server/books"
	"go-app/pkg/auth"
	"go-app/pkg/usr"
)

func main() {
	const SECRET = "secret"
	const ADDRESS = ":1323"

	usrSvc := &UserServiceMemory{}

	s := server.New(nil, SECRET)

	authSvc := auth.New(usrSvc, SECRET)

	routers := []server.RouterGroup{
		iauth.NewAuthGroup(authSvc),
		books.NewBookGroup(SECRET),
	}

	for _, r := range routers {
		s.AddGroup(r)
	}

	s.Start(ADDRESS)
}

type UserServiceMemory struct {
}

func (u *UserServiceMemory) Find(id, pw string) (usr.User, error) {

	if id != "user1" || pw != "1234" {
		return usr.User{}, errors.New("user not found")
	}

	return usr.User{
		Id:    "user1",
		Email: "user1@naver.com",
	}, nil
}

func (u *UserServiceMemory) IsValid(id, pw string) bool {
	return true
}
