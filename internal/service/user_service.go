package service

import (
	"errors"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService implements domain.UserService
type UserService struct {
	repo   domain.UserRepository
	logger *logger.Logger
}

// NewUserService creates a new UserService
func NewUserService(repo domain.UserRepository, logger *logger.Logger) domain.UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get user by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(username string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// List retrieves a list of users with pagination
func (s *UserService) List(limit, offset int32) ([]*domain.User, error) {
	users, err := s.repo.List(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// Create creates a new user
func (s *UserService) Create(user *domain.User, password string) (*domain.User, error) {
	// Check if username already exists
	existingUser, err := s.repo.GetByUsername(user.Username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		s.logger.Error("Error checking username existence", zap.String("username", user.Username), zap.Error(err))
		return nil, err
	}

	// Check if email already exists
	existingUser, err = s.repo.GetByEmail(user.Email)
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
	createdUser, err := s.repo.Create(user)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return createdUser, nil
}

// Update updates an existing user
func (s *UserService) Update(user *domain.User) (*domain.User, error) {
	// Check if user exists
	existingUser, err := s.repo.GetByID(user.ID)
	if err != nil {
		s.logger.Error("Failed to get user by ID", zap.Int64("id", user.ID), zap.Error(err))
		return nil, err
	}

	// Check if username is being changed and if it already exists
	if user.Username != existingUser.Username {
		userByUsername, err := s.repo.GetByUsername(user.Username)
		if err == nil && userByUsername != nil {
			return nil, domain.ErrUserAlreadyExists
		}
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			s.logger.Error("Error checking username existence", zap.String("username", user.Username), zap.Error(err))
			return nil, err
		}
	}

	// Check if email is being changed and if it already exists
	if user.Email != existingUser.Email {
		userByEmail, err := s.repo.GetByEmail(user.Email)
		if err == nil && userByEmail != nil {
			return nil, domain.ErrUserAlreadyExists
		}
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			s.logger.Error("Error checking email existence", zap.String("email", user.Email), zap.Error(err))
			return nil, err
		}
	}

	// Preserve password hash
	user.PasswordHash = existingUser.PasswordHash

	// Update user
	updatedUser, err := s.repo.Update(user)
	if err != nil {
		s.logger.Error("Failed to update user", zap.Int64("id", user.ID), zap.Error(err))
		return nil, err
	}

	return updatedUser, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id int64, currentPassword, newPassword string) error {
	// Get user
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get user by ID", zap.Int64("id", id), zap.Error(err))
		return err
	}

	// Verify current password
	if !verifyPassword(currentPassword, user.PasswordHash) {
		return domain.ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return err
	}

	// Update password
	err = s.repo.UpdatePassword(id, hashedPassword)
	if err != nil {
		s.logger.Error("Failed to update password", zap.Int64("id", id), zap.Error(err))
		return err
	}

	return nil
}

// Delete deletes a user
func (s *UserService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete user", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}

// ValidateCredentials validates a user's credentials
func (s *UserService) ValidateCredentials(username, password string) (*domain.User, error) {
	// Get user by username
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	// Verify password
	if !verifyPassword(password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

// Helper functions

// hashPassword hashes a password
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against a hash
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
