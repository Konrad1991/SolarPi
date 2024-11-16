package main

import (
	"SolarPi/internal/Server"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	err := Server.StartServer(":8080", "Database.db")
	if err != nil {
		fmt.Println(err)
	}
}
