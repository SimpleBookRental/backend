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
// RegisterUser godoc
// @Summary      Register user
// @Description  Register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreate  true  "User register payload"
// @Success      201   {object}  models.User
// @Failure      400   {object}  models.ErrorResponse
// @Router       /api/v1/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var userCreate models.UserCreate
	if err := ctx.ShouldBindJSON(&userCreate); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	user, err := c.userService.Register(&userCreate)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Created(ctx, user)
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

	utils.OK(ctx, user)
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
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.OK(ctx, users)
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
		utils.BadRequest(ctx, err.Error())
		return
	}
	// Enforce email change restrictions
	ctxUserID, ctxUserIDExists := ctx.Get("user_id")
	ctxRole, _ := ctx.Get("role")
	if userUpdate.Email != "" {
		// Regular users cannot change their own email
		if ctxRole == models.UserRole {
			utils.Forbidden(ctx, "You are not allowed to change your email")
			return
		}
		// Admin cannot change their own email
		if ctxRole == models.AdminRole && ctxUserIDExists && ctxUserID.(string) == id {
			utils.Forbidden(ctx, "Admin cannot change their own email")
			return
		}
	}

	user, err := c.userService.Update(id, &userUpdate)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, user)
}

//
// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]bool
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /api/v1/users/{id} [delete]
// @Security     BearerAuth
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.userService.Delete(id)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, gin.H{"deleted": true})
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
		utils.BadRequest(ctx, err.Error())
		return
	}

	response, err := c.userService.Login(&userLogin, c.tokenRepo)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, response)
}
