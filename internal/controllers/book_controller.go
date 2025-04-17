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

//
// CreateBook godoc
// @Summary      Create book
// @Description  Create a new book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        book  body      models.BookCreate  true  "Book create payload"
// @Success      201   {object}  models.Book
// @Failure      400   {object}  models.ErrorResponse
// @Router       /api/v1/books [post]
// @Security     BearerAuth
func (c *BookController) Create(ctx *gin.Context) {
	var bookCreate models.BookCreate
	if err := ctx.ShouldBindJSON(&bookCreate); err != nil {
		utils.BadRequest(ctx, err.Error())
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
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Created(ctx, book)
}

//
// GetBookByID godoc
// @Summary      Get book by ID
// @Description  Get a book by its ID
// @Tags         Books
// @Produce      json
// @Param        id   path      string  true  "Book ID"
// @Success      200  {object}  models.Book
// @Failure      404  {object}  models.ErrorResponse
// @Router       /api/v1/books/{id} [get]
// @Security     BearerAuth
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

	utils.OK(ctx, book)
}

//
// GetAllBooks godoc
// @Summary      Get all books
// @Description  Get all books
// @Tags         Books
// @Produce      json
// @Success      200  {array}   models.Book
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/books [get]
// @Security     BearerAuth
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
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.OK(ctx, books)
}

//
// UpdateBook godoc
// @Summary      Update book
// @Description  Update a book by ID
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Book ID"
// @Param        book  body      models.BookUpdate   true  "Book update payload"
// @Success      200   {object}  models.Book
// @Failure      400   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /api/v1/books/{id} [put]
// @Security     BearerAuth
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
		utils.BadRequest(ctx, err.Error())
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
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, book)
}

//
// DeleteBook godoc
// @Summary      Delete book
// @Description  Delete a book by ID
// @Tags         Books
// @Produce      json
// @Param        id   path      string  true  "Book ID"
// @Success      200  {object}  map[string]bool
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /api/v1/books/{id} [delete]
// @Security     BearerAuth
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
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, gin.H{"deleted": true})
}
