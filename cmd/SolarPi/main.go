package main

import (
    "fmt"
    "time"
    "github.com/corrupt/go-smbus"
)

func main() {
    addr := 0x38
    bus := 0x01
    smb, err := smbus.New(addr, bus)
    if err != nil {
        panic("Could not create smbus")
    }
    defer smb.Bus_close()

    // Initialise
    time.Sleep(0.5 * time.Second) 
    smb.Read_byte_data(0x712)

    fmt.Printf("Hello World")
}
