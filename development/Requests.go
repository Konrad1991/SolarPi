package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Replace with the path to your test file
	filePath := "./development/testfile.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	fileField, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		fmt.Printf("Failed to create form file: %v\n", err)
		return
	}
	_, err = io.Copy(fileField, file)
	if err != nil {
		fmt.Printf("Failed to copy file content: %v\n", err)
		return
	}
	_ = writer.WriteField("FileName", "TestTest.txt")

	writer.Close()

	// Make the POST request
	req, err := http.NewRequest("POST", "http://localhost:8080/UploadFile", &requestBody)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)
}
