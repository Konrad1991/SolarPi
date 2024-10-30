// NOTE: https://github.com/mattn/go-sqlite3/issues/803
// export CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
package Server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type File struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	FileName  string    `json:"file_name"`
	Size      int64     `json:"file_size"`
	Date      time.Time `json:"date"`
	Extension string    `json:"extension"`
}

var DB *gorm.DB

// Init database
func initDB() error {
	DB, err := gorm.Open("sqlite3", "./internal/Server/test.db")
	if err != nil {
		return errors.New("Failed to connect to database")
	}
	DB.AutoMigrate(&File{})
	return nil
}

// create gin router
func createRouter(ip string) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	// r.SetTrustedProxies([]string{"192.168.1.2"}
	return r
}

// Create routes
func createRoutes(r *gin.Engine) {
	r.GET("/GetFiles", getFiles)
	r.GET("/GetFile/:id", getFile)
	r.POST("/CreateFile", createFile) // NOTE: an empty file is created
	r.PUT("/UpdateFile/:id", updateFile)
	r.DELETE("/DeleteFile/:id", deleteFile)
	r.POST("/UploadFile", uploadFile)
}

func Server(ip_addr string) (error, *gin.Engine) {
	initDB()
	r := createRouter(ip_addr)
	createRoutes(r)
	return nil, r
}

// Definition of routes
// ===============================================================================
func uploadFile(c *gin.Context) {
	fmt.Println("uploadFile0")
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse multipart form"})
		return
	}
	fmt.Println("uploadFile1")
	file_name := c.Request.FormValue("FileName")
	file, _, err := c.Request.FormFile("file")

	fmt.Println("uploadFile2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in request"})
		return
	}

	fmt.Println("uploadFile3")
	fmt.Println("file_name: ", file_name)
	defer file.Close()
	file_struct := File{
		FileName:  file_name,
		Size:      0, // TODO: get file size
		Date:      time.Now(),
		Extension: "TODO: get file extension",
	}

	fmt.Println(DB == nil)
	DB.Create(&file_struct)
	// TODO: where is the file saved? saved to DB?
	c.JSON(http.StatusCreated, file_struct)
}

func getFiles(c *gin.Context) { // NOTE: request all files
	var files []File
	DB.Find(&files)
	c.JSON(http.StatusOK, files)
}

func getFile(c *gin.Context) { // NOTE: request one file
	id := c.Param("id")
	var file File
	if err := DB.First(&file, id).Error; err != nil {
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
	DB.Create(&file)
	c.JSON(http.StatusCreated, file)
}

func updateFile(c *gin.Context) {
	id := c.Param("id")
	var file File
	if err := DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if err := c.BindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	DB.Save(&file)
	c.JSON(http.StatusOK, file)
}

func deleteFile(c *gin.Context) {
	id := c.Param("id")
	var file File
	if err := DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	DB.Delete(&file)
	c.JSON(http.StatusOK, gin.H{"message": "File could not be deleted"})
}
