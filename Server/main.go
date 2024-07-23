// https://raw.githubusercontent.com/gin-gonic/gin/master/docs/doc.md
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var db *gorm.DB
var err error

func main() {
	r := gin.Default()
  r.SetTrustedProxies(nil)
  // r.SetTrustedProxies([]string{"192.168.1.2"})

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Book{})

	r.GET("/books", getBooks)
	r.GET("/books/:id", getBook)
	r.POST("/books", createBook)
	r.PUT("/books/:id", updateBook)
	r.DELETE("/books/:id", deleteBook)

	r.POST("/books/upload", uploadFile)

	r.Run(":8080")
}

func uploadFile(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse multipart form"})
		return
	}
	title := c.Request.FormValue("title")
	author := c.Request.FormValue("author")
	file, _, err := c.Request.FormFile("file")

fileTest := file.Open()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in request"})
		return
	}
	defer file.Close()
	book := Book{
		Title:  title,
		Author: author,
	}
	db.Create(&book)
	c.JSON(http.StatusCreated, book)
}

func getBooks(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	db.Create(&book)
	c.JSON(http.StatusCreated, book)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	db.Save(&book)
	c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	db.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

