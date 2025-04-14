package services

import (
	"errors"
	"fmt"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/pkg/utils"
)

// UserService handles business logic for users
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Create creates a new user
func (s *UserService) Create(userCreate *models.UserCreate) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(userCreate.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create user
	user := &models.User{
		Name:     userCreate.Name,
		Email:    userCreate.Email,
		Password: userCreate.Password, // In a real application, you would hash the password
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// GetByID gets a user by ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	if !utils.IsValidUUID(id) {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetAll gets all users
func (s *UserService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// Update updates a user
func (s *UserService) Update(id string, userUpdate *models.UserUpdate) (*models.User, error) {
	if !utils.IsValidUUID(id) {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}
	if userUpdate.Email != "" {
		// Check if email already exists
		if userUpdate.Email != user.Email {
			existingUser, err := s.userRepo.FindByEmail(userUpdate.Email)
			if err != nil {
				return nil, fmt.Errorf("error checking existing user: %w", err)
			}
			if existingUser != nil {
				return nil, errors.New("email already exists")
			}
			user.Email = userUpdate.Email
		}
	}
	if userUpdate.Password != "" {
		user.Password = userUpdate.Password // In a real application, you would hash the password
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

// Delete deletes a user
func (s *UserService) Delete(id string) error {
	if !utils.IsValidUUID(id) {
		return errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
