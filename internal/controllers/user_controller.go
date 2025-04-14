package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService services.UserServiceInterface
	tokenRepo   repositories.TokenRepositoryInterface
}

// NewUserController creates a new user controller
func NewUserController(userService services.UserServiceInterface, tokenRepo repositories.TokenRepositoryInterface) *UserController {
	return &UserController{userService: userService, tokenRepo: tokenRepo}
}

// Create handles the creation of a new user
func (c *UserController) Create(ctx *gin.Context) {
	var userCreate models.UserCreate
	if err := ctx.ShouldBindJSON(&userCreate); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	user, err := c.userService.Create(&userCreate)
	if err != nil {
		utils.BadRequest(ctx, "Failed to create user", err.Error())
		return
	}

	utils.Created(ctx, "User created successfully", user)
}

// GetByID handles getting a user by ID
func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.userService.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.OK(ctx, "User retrieved successfully", user)
}

// GetAll handles getting all users
func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.userService.GetAll()
	if err != nil {
		utils.InternalServerError(ctx, "Failed to retrieve users", err.Error())
		return
	}

	utils.OK(ctx, "Users retrieved successfully", users)
}

// Update handles updating a user
func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var userUpdate models.UserUpdate
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	user, err := c.userService.Update(id, &userUpdate)
	if err != nil {
		utils.BadRequest(ctx, "Failed to update user", err.Error())
		return
	}

	utils.OK(ctx, "User updated successfully", user)
}

// Delete handles deleting a user
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.userService.Delete(id)
	if err != nil {
		utils.BadRequest(ctx, "Failed to delete user", err.Error())
		return
	}

	utils.OK(ctx, "User deleted successfully", nil)
}

// Login handles user login
func (c *UserController) Login(ctx *gin.Context) {
	var userLogin models.UserLogin
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	response, err := c.userService.Login(&userLogin, c.tokenRepo)
	if err != nil {
		utils.BadRequest(ctx, "Login failed", err.Error())
		return
	}

	utils.OK(ctx, "Login successful", response)
}
