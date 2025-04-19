package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// BookUserController handles operations that involve both books and users
type BookUserController struct {
	bookUserService services.BookUserServiceInterface
	userRepo        repositories.UserRepositoryInterface
}

// NewBookUserController creates a new book-user controller
func NewBookUserController(bookUserService services.BookUserServiceInterface, userRepo repositories.UserRepositoryInterface) *BookUserController {
	return &BookUserController{
		bookUserService: bookUserService,
		userRepo:        userRepo,
	}
}

//
// TransferBookOwnership godoc
// @Summary      Transfer book ownership
// @Description  Transfer a book from one user to another
// @Tags         Book-User Operations
// @Accept       json
// @Produce      json
// @Param        id      path      string                        true  "Book ID"
// @Param        request body      models.BookTransferRequest    true  "Transfer request"
// @Success      200     {object}  map[string]bool
// @Failure      400     {object}  models.ErrorResponse
// @Failure      403     {object}  models.ErrorResponse
// @Router       /api/v1/books/{id}/transfer [post]
// @Security     BearerAuth
func (c *BookUserController) TransferBookOwnership(ctx *gin.Context) {
	bookID := ctx.Param("id")
	
	var request struct {
		FromUserID string `json:"from_user_id" binding:"required,uuid"`
		ToUserID   string `json:"to_user_id" binding:"required,uuid"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	
	// Get user role from context
	role, roleExists := ctx.Get("role")
	
	// Only admin can transfer books between any users
	// Regular users can only transfer books they own
if roleExists && role == models.UserRole {
	// Get user ID from context
	userID, userIDExists := ctx.Get("user_id")
	if userIDExists && userID.(string) != request.FromUserID {
		utils.Forbidden(ctx, "You can only transfer books that you own")
		return
	}
	// Ensure target user is a regular user
	targetUser, err := c.userRepo.FindByID(request.ToUserID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	if targetUser == nil {
		utils.BadRequest(ctx, "target user not found")
		return
	}
	if targetUser.Role != models.UserRole {
		utils.Forbidden(ctx, "You can only transfer books to other users")
		return
	}
}
	
	// Call service to transfer book ownership
	err := c.bookUserService.TransferBookOwnership(bookID, request.FromUserID, request.ToUserID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	
	utils.OK(ctx, gin.H{"transferred": true})
}

//
// CreateBookWithUser godoc
// @Summary      Create book with user
// @Description  Create a book and associate it with a user
// @Tags         Book-User Operations
// @Accept       json
// @Produce      json
// @Param        request body      models.BookCreateRequest  true  "Book create with user payload"
// @Success      201     {object}  models.Book
// @Failure      400     {object}  models.ErrorResponse
// @Failure      403     {object}  models.ErrorResponse
// @Router       /api/v1/book-users [post]
// @Security     BearerAuth
func (c *BookUserController) CreateBookWithUser(ctx *gin.Context) {
	var request struct {
		Title       string `json:"title" binding:"required"`
		Author      string `json:"author" binding:"required"`
		ISBN        string `json:"isbn" binding:"required"`
		Description string `json:"description"`
		UserID      string `json:"user_id" binding:"required,uuid"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	
	// Get user role from context
	role, roleExists := ctx.Get("role")
	
	// Only admin can create books for any user
	// Regular users can only create books for themselves
	if roleExists && role == models.UserRole {
		// Get user ID from context
		userID, userIDExists := ctx.Get("user_id")
		if userIDExists && userID.(string) != request.UserID {
			utils.Forbidden(ctx, "You can only create books for yourself")
			return
		}
	}
	
	// Create book request
	bookCreate := &models.BookCreateRequest{
		Title:       request.Title,
		Author:      request.Author,
		ISBN:        request.ISBN,
		Description: request.Description,
	}
	
	// Call service to create book with user
	book, err := c.bookUserService.CreateBookWithUser(bookCreate, request.UserID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	
	utils.Created(ctx, book)
}
