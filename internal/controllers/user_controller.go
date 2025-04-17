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

//
// CreateUser godoc
// @Summary      Create user
// @Description  Create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreate  true  "User create payload"
// @Success      201   {object}  models.User
// @Failure      400   {object}  models.ErrorResponse
// @Router       /api/v1/users [post]
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

//
// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get a user by their ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  models.ErrorResponse
// @Router       /api/v1/users/{id} [get]
func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.userService.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.OK(ctx, "User retrieved successfully", user)
}

//
// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get all users
// @Tags         Users
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/users [get]
// @Security     BearerAuth
func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.userService.GetAll()
	if err != nil {
		utils.InternalServerError(ctx, "Failed to retrieve users", err.Error())
		return
	}

	utils.OK(ctx, "Users retrieved successfully", users)
}

//
// UpdateUser godoc
// @Summary      Update user
// @Description  Update a user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        user  body      models.UserUpdate   true  "User update payload"
// @Success      200   {object}  models.User
// @Failure      400   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /api/v1/users/{id} [put]
// @Security     BearerAuth
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

//
// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /api/v1/users/{id} [delete]
// @Security     BearerAuth
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.userService.Delete(id)
	if err != nil {
		utils.BadRequest(ctx, "Failed to delete user", err.Error())
		return
	}

	utils.OK(ctx, "User deleted successfully", nil)
}

//
// LoginUser godoc
// @Summary      Login
// @Description  User login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.UserLogin  true  "User login payload"
// @Success      200         {object}  models.LoginResponse
// @Failure      400         {object}  models.ErrorResponse
// @Router       /api/v1/login [post]
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
