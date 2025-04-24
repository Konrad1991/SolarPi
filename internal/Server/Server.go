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
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	Name              string    `json:"name" gorm:"unique"`
	Password          string    `json:"password"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UserRootDirectory string    `json:"user_root_directory"`
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

// used for testing
func ResetTestDB() {
	if database != nil {
		database.DropTableIfExists(&User{}, &File{})
		database.AutoMigrate(&User{}, &File{})
	}
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
	// TODO: check whether https.Server is existing
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
	// r.PUT("/UpdateUser/:id", updateUser)
	r.DELETE("/DeleteUser/:Name", deleteUser)
	r.GET("/GetAllUsers", getAllUsers)
	r.POST("/Login", loginUser)
	r.DELETE("/DeleteUserByID/:id", deleteUserByID)
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

	// TODO: Later in production use http and set nginx as reverse proxy
	// Or get a valid certificate
	go func() {
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
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
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func checkUser(c *gin.Context, user *User) error {
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name is empty"})
		return fmt.Errorf("User name is empty")
	}
	var names []string
	if err := database.Model(&User{}).Select("name").Pluck("name", &names).Error; err != nil {
		return err
	}
	for i := 0; i < len(names); i++ {
		if user.Name == names[i] {
			return fmt.Errorf("User name not unique")
		}
	}
	return nil
}

func createUser(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse multipart form"})
		return
	}
	name := c.Request.FormValue("Name")
	password := c.Request.FormValue("Password")
	password_hashed, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hashing failed"})
		return
	}
	t := time.Now()
	user_root_dir := c.Request.FormValue("UserRootDirectory")
	new_user := User{
		Name:              name,
		Password: password_hashed,
		CreatedAt:         t,
		UpdatedAt:         t,
		UserRootDirectory: user_root_dir,
	}

	err = checkUser(c, &new_user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = database.Create(&new_user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, new_user)
}

func deleteUser(c *gin.Context) {
	name := c.Param("Name")
	var user User
	err := database.Where("name = ?", name).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	database.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}

func deleteUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := database.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	database.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func getAllUsers(c *gin.Context) {
	var users []User
	database.Find(&users)
	c.JSON(http.StatusOK, users)
}

func loginUser(c *gin.Context) {
	// TODO:
}
