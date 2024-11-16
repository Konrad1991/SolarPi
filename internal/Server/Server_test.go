package Server

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateRouter(t *testing.T) {
	got := createRouter()
	want := gin.Default()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("type of Router object is not gin default router")
	}
}

func setupTestRouter() *gin.Engine {
	r := createRouter()
	createRoutes(r)
	initDatabase("test.db")
	return r
}

func uploadFileTest(r *gin.Engine, fileName string, fileContent []byte) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, _ := writer.CreateFormFile("file", fileName)
	fileField.Write(fileContent)
	writer.Close()
	req, _ := http.NewRequest("POST", "/UploadFile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	return res
}

func TestUploadFileAndCheckExistence(t *testing.T) {
	r := setupTestRouter()

	// Step 1: Upload the file
	uploadRes := uploadFileTest(r, "testfile.txt", []byte("testfile content"))

	if uploadRes.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d for file upload, but got %d", http.StatusCreated, uploadRes.Code)
	}

	var fileRecord File
	if err := database.First(&fileRecord, 1).Error; err != nil { // Assuming the ID is 1
		t.Fatalf("Failed to retrieve file from database: %v", err)
	}

	// Check the file content
	if string(fileRecord.Content) != string(fileContent) {
		t.Errorf("Expected file content %q, but got %q", string(fileContent), string(fileRecord.Content))
	}

	// Step 2: Check file existence by ID
	getFileReq, _ := http.NewRequest("GET", "/GetFile/1", nil) // Assuming ID is 1 for the first file
	getFileRes := httptest.NewRecorder()
	r.ServeHTTP(getFileRes, getFileReq)

	if getFileRes.Code != http.StatusOK {
		t.Fatalf("Expected status code %d for file retrieval, but got %d", http.StatusOK, getFileRes.Code)
	}

	// Step 3: Check file list using /GetFiles
	getFilesReq, _ := http.NewRequest("GET", "/GetFiles", nil)
	getFilesRes := httptest.NewRecorder()
	r.ServeHTTP(getFilesRes, getFilesReq)

	if getFilesRes.Code != http.StatusOK {
		t.Fatalf("Expected status code %d for getting files list, but got %d", http.StatusOK, getFilesRes.Code)
	}

	// Optional: Check response body if needed for validation
}
