package repository

import (
	"database/sql"
	"errors"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// CategoryRepository implements domain.CategoryRepository
type CategoryRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(conn *DBConn, logger *logger.Logger) domain.CategoryRepository {
	return &CategoryRepository{
		db:     conn.DB,
		logger: logger,
	}
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id int64) (*domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var category domain.Category
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		r.logger.Error("Failed to get category by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return &category, nil
}

// GetByName retrieves a category by name
func (r *CategoryRepository) GetByName(name string) (*domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE name = $1
	`

	var category domain.Category
	err := r.db.QueryRow(query, name).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		r.logger.Error("Failed to get category by name", zap.String("name", name), zap.Error(err))
		return nil, err
	}

	return &category, nil
}

// List retrieves a list of categories with pagination
func (r *CategoryRepository) List(limit, offset int32) ([]*domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list categories", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan category row", zap.Error(err))
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating category rows", zap.Error(err))
		return nil, err
	}

	return categories, nil
}

// ListAll retrieves all categories
func (r *CategoryRepository) ListAll() ([]*domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to list all categories", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan category row", zap.Error(err))
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating category rows", zap.Error(err))
		return nil, err
	}

	return categories, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(category *domain.Category) (*domain.Category, error) {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		category.Name,
		category.Description,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create category", zap.Error(err))
		return nil, err
	}

	return category, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(category *domain.Category) (*domain.Category, error) {
	query := `
		UPDATE categories
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		category.ID,
		category.Name,
		category.Description,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		r.logger.Error("Failed to update category", zap.Int64("id", category.ID), zap.Error(err))
		return nil, err
	}

	return category, nil
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete category", zap.Int64("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}
