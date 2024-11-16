package main

import (
	"SolarPi/internal/Server"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// TODO: move code from func main in Server.go
// TODO: create gin router with ip address as string.
// TODO: create Server_test.go file
// TODO: write test for create router

func main() {
	err := Server.StartServer(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
