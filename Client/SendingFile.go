package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
  "io"
  "bytes"
)

func main() {
	apiURL := "http://localhost:8080/books/upload"

	filePath := "./Readme.txt" 

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Printf("Error creating form file: %s\n", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file to form: %s\n", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing multipart writer: %s\n", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Printf("Error creating POST request: %s\n", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending POST request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Failed to upload file. Status: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("File uploaded successfully!")
}


