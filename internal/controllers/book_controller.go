package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// BookController handles HTTP requests for books
type BookController struct {
	bookService services.BookServiceInterface
}

// NewBookController creates a new book controller
func NewBookController(bookService services.BookServiceInterface) *BookController {
	return &BookController{bookService: bookService}
}

// Create handles the creation of a new book
func (c *BookController) Create(ctx *gin.Context) {
	var bookCreate models.BookCreate
	if err := ctx.ShouldBindJSON(&bookCreate); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	// Get user ID and role from context
	userID, userIDExists := ctx.Get("user_id")
	role, roleExists := ctx.Get("role")

	// Apply business logic based on role
	if roleExists && userIDExists {
		if role == models.UserRole {
			// If role is USER, always use user_id from context
			bookCreate.UserID = userID.(string)
		} else if role == models.AdminRole {
			if bookCreate.UserID == "" {
				// If role is ADMIN and user_id is not provided, use user_id from context
				bookCreate.UserID = userID.(string)
			}
			// If role is ADMIN and user_id is provided, use user_id from body
		}
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

	// Get user ID and role from context
	userID, userIDExists := ctx.Get("user_id")
	role, roleExists := ctx.Get("role")

	// Apply business logic based on role
	if roleExists && userIDExists {
		if role == models.UserRole {
			// If role is USER, can only view their own books
			if book.UserID != userID.(string) {
				utils.Forbidden(ctx, "You do not have permission to view this book")
				return
			}
		}
		// If role is ADMIN, can view any book
	}

	utils.OK(ctx, "Book retrieved successfully", book)
}

// GetAll handles getting all books
func (c *BookController) GetAll(ctx *gin.Context) {
	// Get user ID and role from context
	userID, userIDExists := ctx.Get("user_id")
	role, roleExists := ctx.Get("role")

	var books []models.Book
	var err error

	// Apply business logic based on role
	if roleExists && userIDExists {
		if role == models.AdminRole {
			// If role is ADMIN, get all books
			books, err = c.bookService.GetAll()
		} else if role == models.UserRole {
			// If role is USER, only get books of that user
			books, err = c.bookService.GetByUserID(userID.(string))
		}
	} else {
		// Fallback for tests or when context is not available
		books, err = c.bookService.GetAll()
	}

	if err != nil {
		utils.InternalServerError(ctx, "Failed to retrieve books", err.Error())
		return
	}

	utils.OK(ctx, "Books retrieved successfully", books)
}

// Update handles updating a book
func (c *BookController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get the book first to check ownership
	existingBook, err := c.bookService.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	var bookUpdate models.BookUpdate
	if err := ctx.ShouldBindJSON(&bookUpdate); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	// Get user ID and role from context
	userID, userIDExists := ctx.Get("user_id")
	role, roleExists := ctx.Get("role")

	// Apply business logic based on role
	if roleExists && userIDExists {
		if role == models.UserRole {
			// If role is USER, can only update their own books
			if existingBook.UserID != userID.(string) {
				utils.Forbidden(ctx, "You do not have permission to update this book")
				return
			}

			// Do not allow changing user_id
			bookUpdate.UserID = existingBook.UserID
		}
		// If role is ADMIN, can update any book
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

	// Get the book first to check ownership
	existingBook, err := c.bookService.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	// Get user ID and role from context
	userID, userIDExists := ctx.Get("user_id")
	role, roleExists := ctx.Get("role")

	// Apply business logic based on role
	if roleExists && userIDExists {
		if role == models.UserRole {
			// If role is USER, can only delete their own books
			if existingBook.UserID != userID.(string) {
				utils.Forbidden(ctx, "You do not have permission to delete this book")
				return
			}
		}
		// If role is ADMIN, can delete any book
	}

	err = c.bookService.Delete(id)
	if err != nil {
		utils.BadRequest(ctx, "Failed to delete book", err.Error())
		return
	}

	utils.OK(ctx, "Book deleted successfully", nil)
}
