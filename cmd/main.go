package main

import (
	"SolarPi/internal/Server"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// TODO: move code from func main in Server.go
// TODO: create gin router with ip address as string.
//        --> check within this function that string is a valid ip address
// TODO: create Server_test.go file
// TODO: write test for create router

func main() {
	err, router := Server.Server("IP address")
	Server.DB, err = gorm.Open("sqlite3", "./internal/Server/test.db")
	if err != nil {
		panic("Failed to connect to database")
	}

	Server.DB.AutoMigrate(&Server.File{})
	defer Server.DB.Close()
	if err != nil {
		panic(err)
	}
	router.Run(":8080")
	for {
	}
}
