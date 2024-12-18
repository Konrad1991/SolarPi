// NOTE: https://github.com/mattn/go-sqlite3/issues/803
// export CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
package Server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	ID                  uint      `json:"id" gorm:"primary_key"`
	name                string    `json:"name"`
	password_hash       string    `json:"password_hash"`
	public_key          string    `json:"public_key"`
	created_at          time.Time `json:"created_at"`
	updated_at          time.Time `json:"updated_at"`
	user_root_directory string    `json:"user_root_directory"`
}

type File struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	FileName  string    `json:"file_name"`
	Size      int64     `json:"file_size"`
	Date      time.Time `json:"date"`
	Extension string    `json:"extension"`
	Content   []byte    `json:"content"`
}

var database *gorm.DB

// Init database
func initDatabase(databaseName string) error {
	var err error
	database, err = gorm.Open("sqlite3", databaseName)
	if err != nil {
		return errors.New("Failed to connect to database")
	}
	database.AutoMigrate(&File{})
	database.AutoMigrate(&User{})
	return nil
}

// Create gin router
func createRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome SolarPi Server")
	})
	return router
}

// Create server
func createServer(ip string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    ip, //":8080",
		Handler: router,
	}
}

// Create routes
func createRoutes(r *gin.Engine) {
	r.GET("/GetFiles", getFiles)
	r.GET("/GetFile/:id", getFile)
	r.POST("/CreateFile", createFile) // NOTE: an empty file is created
	r.PUT("/UpdateFile/:id", updateFile)
	r.DELETE("/DeleteFile/:id", deleteFile)
	r.POST("/UploadFile", uploadFile)

	r.POST("/CreateUser", createUser)
	r.PUT("/UpdateUser/:id", updateUser)
	r.DELETE("/DeleteUser/:id", deleteUser)
	r.GET("/GetAllUsers", getAllUsers)
}

func StartServer(ip_addr string, databaseName string) error {
	// Database
	if err := initDatabase(databaseName); err != nil {
		return err
	}
	defer func() {
		sqlDB := database.DB()
		sqlDB.Close()
	}()

	// Router & Server
	r := createRouter()
	createRoutes(r)
	srv := createServer(ip_addr, r)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %w", err)
	}
	return nil
}

// Definition of routes
// ===============================================================================

// File routes
// ===============================================================================
func uploadFile(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse multipart form"})
		return
	}
	file_name := c.Request.FormValue("FileName")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in request"})
		return
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}
	splitted_file := strings.SplitAfter(header.Filename, ".")
	ext := "none"
	if len(splitted_file) != 1 {
		ext = splitted_file[len(splitted_file)-1]
	}
	file_struct := File{
		FileName:  file_name,
		Size:      int64(len(fileContent)),
		Date:      time.Now(),
		Extension: ext,
		Content:   fileContent,
	}
	database.Create(&file_struct)
	c.JSON(http.StatusCreated, file_struct)
}

func getFiles(c *gin.Context) { // NOTE: request all files
	var files []File
	database.Find(&files)
	c.JSON(http.StatusOK, files)
}

func getFile(c *gin.Context) { // NOTE: request one file
	id := c.Param("id")
	var file File
	if err := database.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}

func createFile(c *gin.Context) { // NOTE: create a new empty file
	var file File
	if err := c.BindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	database.Create(&file)
	c.JSON(http.StatusCreated, file)
}

func updateFile(c *gin.Context) {
	id := c.Param("id")
	var file File
	if err := database.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if err := c.BindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	database.Save(&file)
	c.JSON(http.StatusOK, file)
}

func deleteFile(c *gin.Context) {
	id := c.Param("id")
	var file File
	if err := database.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	database.Delete(&file)
	c.JSON(http.StatusOK, gin.H{"message": "File could not be deleted"})
}

// User routes
// ===============================================================================
func createUser(c *gin.Context) {
}
