package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// BookController handles HTTP requests for books
type BookController struct {
	bookService *services.BookService
}

// NewBookController creates a new book controller
func NewBookController(bookService *services.BookService) *BookController {
	return &BookController{bookService: bookService}
}

// Create handles the creation of a new book
func (c *BookController) Create(ctx *gin.Context) {
	var bookCreate models.BookCreate
	if err := ctx.ShouldBindJSON(&bookCreate); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	book, err := c.bookService.Create(&bookCreate)
	if err != nil {
		utils.BadRequest(ctx, "Failed to create book", err.Error())
		return
	}

	utils.Created(ctx, "Book created successfully", book)
}

// GetByID handles getting a book by ID
func (c *BookController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	book, err := c.bookService.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.OK(ctx, "Book retrieved successfully", book)
}

// GetAll handles getting all books
func (c *BookController) GetAll(ctx *gin.Context) {
	books, err := c.bookService.GetAll()
	if err != nil {
		utils.InternalServerError(ctx, "Failed to retrieve books", err.Error())
		return
	}

	utils.OK(ctx, "Books retrieved successfully", books)
}

// GetByUserID handles getting all books by user ID
func (c *BookController) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	books, err := c.bookService.GetByUserID(userID)
	if err != nil {
		utils.BadRequest(ctx, "Failed to retrieve books", err.Error())
		return
	}

	utils.OK(ctx, "Books retrieved successfully", books)
}

// Update handles updating a book
func (c *BookController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var bookUpdate models.BookUpdate
	if err := ctx.ShouldBindJSON(&bookUpdate); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	book, err := c.bookService.Update(id, &bookUpdate)
	if err != nil {
		utils.BadRequest(ctx, "Failed to update book", err.Error())
		return
	}

	utils.OK(ctx, "Book updated successfully", book)
}

// Delete handles deleting a book
func (c *BookController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.bookService.Delete(id)
	if err != nil {
		utils.BadRequest(ctx, "Failed to delete book", err.Error())
		return
	}

	utils.OK(ctx, "Book deleted successfully", nil)
}
