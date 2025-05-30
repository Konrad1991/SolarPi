package Server

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Test router creation
// ====================================================================
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
	ResetTestDB()
	return r
}

// Test User code
// ====================================================================
func createUserTest(r *gin.Engine, name, passwordHash, userRootDir string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("Name", name)
	writer.WriteField("Password", passwordHash)
	writer.WriteField("UserRootDirectory", userRootDir)
	writer.Close()
	req, _ := http.NewRequest("POST", "/CreateUser", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	return res
}

func TestCreateUserAndVerify(t *testing.T) {
	r := setupTestRouter()
	name := "testuser"
	password := "password"
	userRootDir := "/home/testuser"
	createRes := createUserTest(r, name, password, userRootDir)
	if createRes.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d for user creation, but got %d",
			http.StatusCreated, createRes.Code)
	}
	var userRecord User
	if err := database.First(&userRecord, "Name = ?", name).Error; err != nil {
		t.Fatalf("Failed to retrieve user from database: %v", err)
	}
	if userRecord.Name != name {
		t.Errorf("Expected name %s, got %s", name, userRecord.Name)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRecord.Password), []byte(password)); err != nil {
		t.Errorf("Password hash does not match: %v", err)
	}

	if userRecord.UserRootDirectory != userRootDir {
		t.Errorf("Expected root dir %s, got %s", userRootDir, userRecord.UserRootDirectory)
	}
}

func TestDeleteUser(t *testing.T) {
	r := setupTestRouter()
	name := "userToDelete"
	passwordHash := "hashedpassword"
	userRootDir := "/home/userToDelete"
	createUserTest(r, name, passwordHash, userRootDir)
	req, _ := http.NewRequest("DELETE", "/DeleteUser/"+name, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("Expected status code %d for user deletion, but got %d", http.StatusOK, res.Code)
	}
	var userRecord User
	if err := database.First(&userRecord, "name = ?", name).Error; err == nil {
		t.Errorf("User with name %s was not deleted", name)
	}
}

func TestGetAllUsers(t *testing.T) {
	r := setupTestRouter()
	users := []struct {
		Name         string
		Password string
		UserRootDir  string
	}{
		{"user1", "hash1", "/home/user1"},
		{"user2", "hash2", "/home/user2"},
		{"user3", "hash3", "/home/user3"},
	}
	for _, u := range users {
		createUserTest(r, u.Name, u.Password, u.UserRootDir)
	}
	req, _ := http.NewRequest("GET", "/GetAllUsers", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("Expected status code %d for getting all users, but got %d", http.StatusOK, res.Code)
	}
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
	uploadRes := uploadFileTest(r, "testfile.txt", []byte("testfile content"))
	if uploadRes.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d for file upload, but got %d", http.StatusCreated, uploadRes.Code)
	}
	var fileRecord File
	if err := database.First(&fileRecord, 1).Error; err != nil { // Assuming the ID is 1
		t.Fatalf("Failed to retrieve file from database: %v", err)
	}

 	// Check the file content
 	fileContent := []byte("testfile content")
 	if string(fileRecord.Content) != string(fileContent) {
 		t.Errorf("Expected file content %q, but got %q", string(fileContent), string(fileRecord.Content))
 	}
 	// Check file existens by id
 	getFileReq, _ := http.NewRequest("GET", "/GetFile/1", nil)
 	getFileRes := httptest.NewRecorder()
 	r.ServeHTTP(getFileRes, getFileReq)

 	if getFileRes.Code != http.StatusOK {
 		t.Fatalf("Expected status code %d for file retrieval, but got %d", http.StatusOK, getFileRes.Code)
 	}
 	// Check file list using /GetFiles
 	getFilesReq, _ := http.NewRequest("GET", "/GetFiles", nil)
 	getFilesRes := httptest.NewRecorder()
 	r.ServeHTTP(getFilesRes, getFilesReq)
 	if getFilesRes.Code != http.StatusOK {
 		t.Fatalf("Expected status code %d for getting files list, but got %d", http.StatusOK, getFilesRes.Code)
 	}
 }
