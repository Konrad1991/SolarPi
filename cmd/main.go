package main

import (
  "SolarPi/internal/Server"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
  err, router:= Server.Server("IP address")
  Server.DB, err = gorm.Open("sqlite3", "./internal/Server/test.db")
	if err != nil {
	  panic("Failed to connect to database")
	}

  Server.DB.AutoMigrate(&Server.File{})
  defer Server.DB.Close()
  if (err != nil) {
    panic(err)
  }
  router.Run(":8080")
  for {

  }
}
