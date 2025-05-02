package api

import (
	"net/http"

	"github.com/SimpleBookRental/backend/internal/service"
	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService service.AuthService
	jwtService  *auth.JWTService
	logger      *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService, jwtService *auth.JWTService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email     string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Register a new user with member role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "User registration information"
// @Success      201   {object}  domain.User
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      409   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      domain.RoleMember,
	}

	createdUser, err := h.authService.Register(user, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		SendError(c, err)
		return
	}

	SendCreated(c, createdUser, "User registered successfully")
}

// Login handles user login
// @Summary      Login user
// @Description  Authenticate user and return access & refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "User credentials"
// @Success      200          {object}  TokenResponse
// @Failure      400          {object}  domain.ErrorResponse
// @Failure      401          {object}  domain.ErrorResponse
// @Failure      500          {object}  domain.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		SendError(c, err)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	SendSuccess(c, response, "Login successful")
}

// RefreshToken handles token refresh
// @Summary      Refresh token
// @Description  Generate new access and refresh tokens using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  body      RefreshTokenRequest  true  "Refresh token"
// @Success      200    {object}  TokenResponse
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      401    {object}  domain.ErrorResponse
// @Failure      500    {object}  domain.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh token", zap.Error(err))
		SendError(c, err)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	SendSuccess(c, response, "Token refreshed successfully")
}

// Logout handles user logout
// @Summary      Logout user
// @Description  Invalidate the current user's token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	// Extract token from header
	token := authHeader[7:] // Remove "Bearer " prefix

	err := h.authService.Logout(token)
	if err != nil {
		h.logger.Error("Failed to logout", zap.Error(err))
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
