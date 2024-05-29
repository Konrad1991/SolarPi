package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)


type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {
	apiURL := "http://localhost:8080/books"

	newBook := Book{
		Title:  "New Book1",
		Author: "Jane Doe1",
	}

	jsonBody, err := json.Marshal(newBook)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Failed to create book. Status: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("Book created successfully!")

}

