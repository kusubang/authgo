package main

import (
	"fmt"
	"go-app/pkg/auth"
	"go-app/pkg/usr"
	"log"
)

func main() {

	auth := auth.New(
		&usr.UserServiceMem{},
	)
	accessToken, refreshToken, err := auth.Login("user1", "user1@naver.com", "1234")

	if err != nil {
		log.Fatal(err)
		return
	}

	c1, err := auth.IsValid(accessToken + "-")
	if err != nil {
		log.Printf("invalid access token\n")
	}
	c2, err := auth.IsValid(refreshToken)
	if err != nil {
		log.Printf("invalid refresh token\n")
	}

	fmt.Println(c1)
	fmt.Println(c2)
	// fmt.Println(accessToken)
	// fmt.Println(refreshToken)
}
