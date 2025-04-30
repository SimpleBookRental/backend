package service

import (
	"errors"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// CategoryService implements domain.CategoryService
type CategoryServiceImpl struct {
	repo   domain.CategoryRepository
	logger *logger.Logger
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(repo domain.CategoryRepository, logger *logger.Logger) domain.CategoryService {
	return &CategoryServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

// GetByID retrieves a category by ID
func (s *CategoryServiceImpl) GetByID(id int64) (*domain.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get category by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return category, nil
}

// GetByName retrieves a category by name
func (s *CategoryServiceImpl) GetByName(name string) (*domain.Category, error) {
	category, err := s.repo.GetByName(name)
	if err != nil {
		s.logger.Error("Failed to get category by name", zap.String("name", name), zap.Error(err))
		return nil, err
	}
	return category, nil
}

// List retrieves a list of categories with pagination
func (s *CategoryServiceImpl) List(limit, offset int32) ([]*domain.Category, error) {
	categories, err := s.repo.List(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list categories", zap.Error(err))
		return nil, err
	}
	return categories, nil
}

// ListAll retrieves all categories
func (s *CategoryServiceImpl) ListAll() ([]*domain.Category, error) {
	categories, err := s.repo.ListAll()
	if err != nil {
		s.logger.Error("Failed to list all categories", zap.Error(err))
		return nil, err
	}
	return categories, nil
}

// Create creates a new category
func (s *CategoryServiceImpl) Create(category *domain.Category) (*domain.Category, error) {
	// Check if category name already exists
	existingCategory, err := s.repo.GetByName(category.Name)
	if err == nil && existingCategory != nil {
		return nil, domain.ErrCategoryAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrCategoryNotFound) {
		s.logger.Error("Error checking category name existence", zap.String("name", category.Name), zap.Error(err))
		return nil, err
	}

	// Create category
	createdCategory, err := s.repo.Create(category)
	if err != nil {
		s.logger.Error("Failed to create category", zap.Error(err))
		return nil, err
	}

	return createdCategory, nil
}

// Update updates an existing category
func (s *CategoryServiceImpl) Update(category *domain.Category) (*domain.Category, error) {
	// Check if category exists
	existingCategory, err := s.repo.GetByID(category.ID)
	if err != nil {
		s.logger.Error("Failed to get category by ID", zap.Int64("id", category.ID), zap.Error(err))
		return nil, err
	}

	// Check if category name is being changed and if it already exists
	if category.Name != existingCategory.Name {
		categoryByName, err := s.repo.GetByName(category.Name)
		if err == nil && categoryByName != nil {
			return nil, domain.ErrCategoryAlreadyExists
		}
		if err != nil && !errors.Is(err, domain.ErrCategoryNotFound) {
			s.logger.Error("Error checking category name existence", zap.String("name", category.Name), zap.Error(err))
			return nil, err
		}
	}

	// Update category
	updatedCategory, err := s.repo.Update(category)
	if err != nil {
		s.logger.Error("Failed to update category", zap.Int64("id", category.ID), zap.Error(err))
		return nil, err
	}

	return updatedCategory, nil
}

// Delete deletes a category
func (s *CategoryServiceImpl) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete category", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}
