package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book model representing a library book
type Book struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	Genre         string `json:"genre" binding:"required"`
	PublishedYear int    `json:"published_year" binding:"required"`
	ISBN          string `json:"isbn" binding:"required"`
	Availability  bool   `json:"availability"`
}

var db *gorm.DB

// Initialize and connect to the database
func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&Book{})
}

// Middleware to handle CORS for frontend integration
func setupCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

// Get all books
func getBooks(c *gin.Context) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// Get a book by ID
func getBookByID(c *gin.Context) {
	id := c.Param("id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Add a new book
func addBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// Update an existing book
func updateBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updateData Book
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update book fields (excluding ID)
	db.Model(&book).Updates(updateData)

	c.JSON(http.StatusOK, book)
}

// Delete a book
func deleteBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func main() {
	initDatabase()

	r := gin.Default()
	r.Use(setupCORS()) // Enable CORS

	// API Routes
	r.GET("/books", getBooks)
	r.GET("/books/:id", getBookByID)
	r.POST("/books", addBook)
	r.PUT("/books/:id", updateBook)
	r.DELETE("/books/:id", deleteBook)

	// Start server
	r.Run(":8080")
}

