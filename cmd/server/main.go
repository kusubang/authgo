package main

import (
	"go-app/internal/server"
)

func main() {
	const SECRET = "secret"
	s := server.New(SECRET)
	s.Start(":1323")
}
