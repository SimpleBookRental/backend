package service

import (
	"errors"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo   domain.UserRepository
	jwtService *auth.JWTService
	logger     *logger.Logger
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo domain.UserRepository, jwtService *auth.JWTService, logger *logger.Logger) AuthService {
	return &AuthServiceImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

// Register registers a new user
func (s *AuthServiceImpl) Register(user *domain.User, password string) (*domain.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(user.Username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		s.logger.Error("Error checking username existence", zap.String("username", user.Username), zap.Error(err))
		return nil, err
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		s.logger.Error("Error checking email existence", zap.String("email", user.Email), zap.Error(err))
		return nil, err
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}
	user.PasswordHash = hashedPassword

	// Set default role if not provided
	if user.Role == "" {
		user.Role = domain.RoleMember
	}

	// Create user
	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return createdUser, nil
}

// Login authenticates a user and returns access and refresh tokens
func (s *AuthServiceImpl) Login(username, password string) (string, string, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", "", domain.ErrInvalidCredentials
		}
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return "", "", err
	}

	// Verify password
	if !verifyPassword(password, user.PasswordHash) {
		return "", "", domain.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return "", "", err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthServiceImpl) RefreshToken(refreshToken string) (string, string, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		s.logger.Error("Invalid refresh token", zap.Error(err))
		return "", "", domain.ErrUnauthorized
	}

	// Check token type
	if claims.TokenType != auth.RefreshToken {
		s.logger.Error("Token is not a refresh token", zap.String("tokenType", string(claims.TokenType)))
		return "", "", domain.ErrUnauthorized
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		s.logger.Error("Failed to get user by ID", zap.Int64("id", claims.UserID), zap.Error(err))
		return "", "", err
	}

	// Generate new tokens
	newAccessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return "", "", err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout logs out a user
func (s *AuthServiceImpl) Logout(token string) error {
	// In a real-world application, you might want to invalidate the token
	// by adding it to a blacklist or using Redis to store invalidated tokens
	// For simplicity, we'll just return nil
	return nil
}
