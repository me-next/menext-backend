package main

import (
	"fmt"
	"github.com/me-next/menext-backend/server"
)

func main() {
	fmt.Println("hello world")
	s := server.New()

	// TODO: maybe handle this error better...
	panic(s.Start(":8080"))
}
