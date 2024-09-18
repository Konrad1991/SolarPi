package main

import (
  "SolarPi/internal/Server"
)

func main() {
  err, router:= Server.Server("IP address")
  db := Server.DB
  defer db.Close()
  if (err != nil) {
    panic(err)
  }
  router.Run(":8080")
  for {

  }
}
