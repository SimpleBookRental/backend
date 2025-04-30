package repository

import (
	"database/sql"
	"errors"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(conn *DBConn, logger *logger.Logger) domain.UserRepository {
	return &UserRepository{
		db:     conn.DB,
		logger: logger,
	}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(limit, offset int32) ([]*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan user row", zap.Error(err))
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating user rows", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, password_hash, first_name, last_name, role, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users
		SET username = $2, email = $3, first_name = $4, last_name = $5, role = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, email, password_hash, first_name, last_name, role, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("Failed to update user", zap.Int64("id", user.ID), zap.Error(err))
		return nil, err
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id, passwordHash)
	if err != nil {
		r.logger.Error("Failed to update user password", zap.Int64("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete user", zap.Int64("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
