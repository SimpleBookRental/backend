package api

import (
	"net/http"
	"strconv"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler handles user requests
type UserHandler struct {
	userService domain.UserService
	jwtService  *auth.JWTService
	logger      *logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService domain.UserService, jwtService *auth.JWTService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Username  string         `json:"username" example:"johndoe"`
	Email     string         `json:"email" example:"john.doe@example.com"`
	FirstName string         `json:"first_name" example:"John"`
	LastName  string         `json:"last_name" example:"Doe"`
	Role      domain.UserRole `json:"role" example:"member"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"oldpassword"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"newpassword123"`
}

// GetByID handles getting a user by ID
// @Summary      Get a user by ID
// @Description  Retrieve a single user by their ID. Users can only view their own profile unless they are admins.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      403  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid user ID"))
		return
	}

	// Check if user is requesting their own profile or is an admin
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if userID.(int64) != id && role != domain.RoleAdmin {
		SendError(c, domain.ErrForbidden)
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, user, "User retrieved successfully")
}

// List handles listing users
// @Summary      List all users
// @Description  Get a paginated list of all users. Only admins can access this endpoint.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.User}
// @Failure      401    {object} domain.ErrorResponse
// @Failure      403    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Security     Bearer
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	// Only admins can list all users
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != string(domain.RoleAdmin) {
		SendError(c, domain.ErrForbidden)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.userService.List(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, users, int64(len(users)), int32(limit), int32(offset), "Users retrieved successfully")
}

// Update handles updating a user
// @Summary      Update a user
// @Description  Update a user's profile information. Users can only update their own profile unless they are admins. Only admins can update user roles.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "User ID"
// @Param        user  body      UpdateUserRequest true  "Updated user information"
// @Success      200   {object}  domain.User
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      401   {object}  domain.ErrorResponse
// @Failure      403   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid user ID"))
		return
	}

	// Check if user is updating their own profile or is an admin
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if userID.(int64) != id && role != domain.RoleAdmin {
		SendError(c, domain.ErrForbidden)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	// Get existing user
	existingUser, err := h.userService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Update user fields
	if req.Username != "" {
		existingUser.Username = req.Username
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.FirstName != "" {
		existingUser.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existingUser.LastName = req.LastName
	}

	// Only admins can update roles
	if req.Role != "" && role == domain.RoleAdmin {
		existingUser.Role = req.Role
	}

	updatedUser, err := h.userService.Update(existingUser)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, updatedUser, "User updated successfully")
}

// Delete handles deleting a user
// @Summary      Delete a user
// @Description  Delete a user from the system. Only admins can delete users.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      403  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid user ID"))
		return
	}

	// Only admins can delete users
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != string(domain.RoleAdmin) {
		SendError(c, domain.ErrForbidden)
		return
	}

	err = h.userService.Delete(id)
	if err != nil {
		h.logger.Error("Failed to delete user", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ChangePassword handles changing a user's password
// @Summary      Change user password
// @Description  Change a user's password. Users can only change their own password.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id         path      int                   true  "User ID"
// @Param        passwords  body      ChangePasswordRequest true  "Password change information"
// @Success      200        {object}  map[string]string
// @Failure      400        {object}  domain.ErrorResponse
// @Failure      401        {object}  domain.ErrorResponse
// @Failure      403        {object}  domain.ErrorResponse
// @Failure      404        {object}  domain.ErrorResponse
// @Failure      500        {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /users/{id}/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid user ID"))
		return
	}

	// Check if user is changing their own password
	userID, exists := c.Get("userID")
	if !exists || userID.(int64) != id {
		SendError(c, domain.ErrForbidden)
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	err = h.userService.ChangePassword(id, req.CurrentPassword, req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to change password", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
