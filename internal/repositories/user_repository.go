package repositories

import (
	"errors"
	"fmt"

	"github.com/SimpleBookRental/backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// GetDB returns the database connection
func (r *UserRepository) GetDB() interface{} {
	return r.db
}

// WithTx returns a new UserRepository with the given transaction
func (r *UserRepository) WithTx(tx interface{}) (UserRepositoryInterface, error) {
	db, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}
	return &UserRepository{db: db}, nil
}

// Ensure UserRepository implements UserRepositoryInterface
var _ UserRepositoryInterface = (*UserRepository)(nil)

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Books").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindAll finds all users
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
