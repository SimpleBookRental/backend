package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// RentalRepository implements domain.RentalRepository
type RentalRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewRentalRepository creates a new RentalRepository
func NewRentalRepository(conn *DBConn, logger *logger.Logger) domain.RentalRepository {
	return &RentalRepository{
		db:     conn.DB,
		logger: logger,
	}
}

// GetByID retrieves a rental by ID
func (r *RentalRepository) GetByID(id int64) (*domain.Rental, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rental_date, r.due_date, r.return_date, r.status,
			   r.created_at, r.updated_at, u.username as user_username, b.title as book_title, b.author as book_author
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		JOIN books b ON r.book_id = b.id
		WHERE r.id = $1
	`

	var rental domain.Rental
	var returnDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&rental.ID,
		&rental.UserID,
		&rental.BookID,
		&rental.RentalDate,
		&rental.DueDate,
		&returnDate,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
		&rental.UserUsername,
		&rental.BookTitle,
		&rental.BookAuthor,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRentalNotFound
		}
		r.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if returnDate.Valid {
		rental.ReturnDate = &returnDate.Time
	}

	return &rental, nil
}

// List retrieves a list of rentals with pagination
func (r *RentalRepository) List(limit, offset int32) ([]*domain.Rental, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rental_date, r.due_date, r.return_date, r.status,
			   r.created_at, r.updated_at, u.username as user_username, b.title as book_title, b.author as book_author
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		JOIN books b ON r.book_id = b.id
		ORDER BY r.rental_date DESC
		LIMIT $1 OFFSET $2
	`

	return r.queryRentals(query, limit, offset)
}

// ListByUser retrieves a list of rentals for a specific user with pagination
func (r *RentalRepository) ListByUser(userID int64, limit, offset int32) ([]*domain.Rental, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rental_date, r.due_date, r.return_date, r.status,
			   r.created_at, r.updated_at, u.username as user_username, b.title as book_title, b.author as book_author
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		JOIN books b ON r.book_id = b.id
		WHERE r.user_id = $1
		ORDER BY r.rental_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list rentals by user", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rentals []*domain.Rental
	for rows.Next() {
		var rental domain.Rental
		var returnDate sql.NullTime

		err := rows.Scan(
			&rental.ID,
			&rental.UserID,
			&rental.BookID,
			&rental.RentalDate,
			&rental.DueDate,
			&returnDate,
			&rental.Status,
			&rental.CreatedAt,
			&rental.UpdatedAt,
			&rental.UserUsername,
			&rental.BookTitle,
			&rental.BookAuthor,
		)
		if err != nil {
			r.logger.Error("Failed to scan rental row", zap.Error(err))
			return nil, err
		}

		if returnDate.Valid {
			rental.ReturnDate = &returnDate.Time
		}

		rentals = append(rentals, &rental)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating rental rows", zap.Error(err))
		return nil, err
	}

	return rentals, nil
}

// ListByBook retrieves a list of rentals for a specific book with pagination
func (r *RentalRepository) ListByBook(bookID int64, limit, offset int32) ([]*domain.Rental, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rental_date, r.due_date, r.return_date, r.status,
			   r.created_at, r.updated_at, u.username as user_username, b.title as book_title, b.author as book_author
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		JOIN books b ON r.book_id = b.id
		WHERE r.book_id = $1
		ORDER BY r.rental_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, bookID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list rentals by book", zap.Int64("bookID", bookID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rentals []*domain.Rental
	for rows.Next() {
		var rental domain.Rental
		var returnDate sql.NullTime

		err := rows.Scan(
			&rental.ID,
			&rental.UserID,
			&rental.BookID,
			&rental.RentalDate,
			&rental.DueDate,
			&returnDate,
			&rental.Status,
			&rental.CreatedAt,
			&rental.UpdatedAt,
			&rental.UserUsername,
			&rental.BookTitle,
			&rental.BookAuthor,
		)
		if err != nil {
			r.logger.Error("Failed to scan rental row", zap.Error(err))
			return nil, err
		}

		if returnDate.Valid {
			rental.ReturnDate = &returnDate.Time
		}

		rentals = append(rentals, &rental)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating rental rows", zap.Error(err))
		return nil, err
	}

	return rentals, nil
}

// ListActive retrieves a list of active rentals with pagination
func (r *RentalRepository) ListActive(limit, offset int32) ([]*domain.Rental, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rental_date, r.due_date, r.return_date, r.status,
			   r.created_at, r.updated_at, u.username as user_username, b.title as book_title, b.author as book_author
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		JOIN books b ON r.book_id = b.id
		WHERE r.status = 'active'
		ORDER BY r.due_date ASC
		LIMIT $1 OFFSET $2
	`

	return r.queryRentals(query, limit, offset)
}

// ListOverdue retrieves a list of overdue rentals with pagination
func (r *RentalRepository) ListOverdue(limit, offset int32) ([]*domain.Rental, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rental_date, r.due_date, r.return_date, r.status,
			   r.created_at, r.updated_at, u.username as user_username, b.title as book_title, b.author as book_author
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		JOIN books b ON r.book_id = b.id
		WHERE r.status = 'active' AND r.due_date < NOW()
		ORDER BY r.due_date ASC
		LIMIT $1 OFFSET $2
	`

	return r.queryRentals(query, limit, offset)
}

// Create creates a new rental
func (r *RentalRepository) Create(rental *domain.Rental) (*domain.Rental, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if book is available
	var availableCopies int32
	err = tx.QueryRow("SELECT available_copies FROM books WHERE id = $1", rental.BookID).Scan(&availableCopies)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		r.logger.Error("Failed to check book availability", zap.Int64("bookID", rental.BookID), zap.Error(err))
		return nil, err
	}

	if availableCopies <= 0 {
		return nil, domain.ErrBookNotAvailable
	}

	// Decrement available copies
	_, err = tx.Exec("UPDATE books SET available_copies = available_copies - 1, updated_at = NOW() WHERE id = $1", rental.BookID)
	if err != nil {
		r.logger.Error("Failed to decrement available copies", zap.Int64("bookID", rental.BookID), zap.Error(err))
		return nil, err
	}

	// Create rental
	query := `
		INSERT INTO rentals (user_id, book_id, rental_date, due_date, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, book_id, rental_date, due_date, return_date, status, created_at, updated_at
	`

	var returnDate sql.NullTime

	err = tx.QueryRow(
		query,
		rental.UserID,
		rental.BookID,
		rental.RentalDate,
		rental.DueDate,
		rental.Status,
	).Scan(
		&rental.ID,
		&rental.UserID,
		&rental.BookID,
		&rental.RentalDate,
		&rental.DueDate,
		&returnDate,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create rental", zap.Error(err))
		return nil, err
	}

	if returnDate.Valid {
		rental.ReturnDate = &returnDate.Time
	}

	// Get user and book details
	err = tx.QueryRow("SELECT username FROM users WHERE id = $1", rental.UserID).Scan(&rental.UserUsername)
	if err != nil {
		r.logger.Error("Failed to get user details", zap.Int64("userID", rental.UserID), zap.Error(err))
		return nil, err
	}

	err = tx.QueryRow("SELECT title, author FROM books WHERE id = $1", rental.BookID).Scan(&rental.BookTitle, &rental.BookAuthor)
	if err != nil {
		r.logger.Error("Failed to get book details", zap.Int64("bookID", rental.BookID), zap.Error(err))
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	return rental, nil
}

// UpdateStatus updates the status of a rental
func (r *RentalRepository) UpdateStatus(id int64, status domain.RentalStatus) (*domain.Rental, error) {
	query := `
		UPDATE rentals
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, book_id, rental_date, due_date, return_date, status, created_at, updated_at
	`

	var rental domain.Rental
	var returnDate sql.NullTime

	err := r.db.QueryRow(query, id, status).Scan(
		&rental.ID,
		&rental.UserID,
		&rental.BookID,
		&rental.RentalDate,
		&rental.DueDate,
		&returnDate,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRentalNotFound
		}
		r.logger.Error("Failed to update rental status", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if returnDate.Valid {
		rental.ReturnDate = &returnDate.Time
	}

	// Get user and book details
	err = r.db.QueryRow("SELECT username FROM users WHERE id = $1", rental.UserID).Scan(&rental.UserUsername)
	if err != nil {
		r.logger.Error("Failed to get user details", zap.Int64("userID", rental.UserID), zap.Error(err))
		return nil, err
	}

	err = r.db.QueryRow("SELECT title, author FROM books WHERE id = $1", rental.BookID).Scan(&rental.BookTitle, &rental.BookAuthor)
	if err != nil {
		r.logger.Error("Failed to get book details", zap.Int64("bookID", rental.BookID), zap.Error(err))
		return nil, err
	}

	return &rental, nil
}

// Return processes the return of a rental
func (r *RentalRepository) Return(id int64) (*domain.Rental, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get rental
	var rental domain.Rental
	var returnDate sql.NullTime
	var bookID int64

	err = tx.QueryRow(`
		SELECT id, user_id, book_id, rental_date, due_date, return_date, status, created_at, updated_at
		FROM rentals
		WHERE id = $1
	`, id).Scan(
		&rental.ID,
		&rental.UserID,
		&bookID,
		&rental.RentalDate,
		&rental.DueDate,
		&returnDate,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRentalNotFound
		}
		r.logger.Error("Failed to get rental", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if rental.Status != domain.RentalStatusActive {
		return nil, domain.ErrRentalNotActive
	}

	// Update rental
	now := time.Now()
	rental.ReturnDate = &now
	rental.Status = domain.RentalStatusReturned

	_, err = tx.Exec(`
		UPDATE rentals
		SET return_date = $2, status = $3, updated_at = NOW()
		WHERE id = $1
	`, id, now, rental.Status)

	if err != nil {
		r.logger.Error("Failed to update rental", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// Increment available copies
	_, err = tx.Exec("UPDATE books SET available_copies = available_copies + 1, updated_at = NOW() WHERE id = $1", bookID)
	if err != nil {
		r.logger.Error("Failed to increment available copies", zap.Int64("bookID", bookID), zap.Error(err))
		return nil, err
	}

	// Get user and book details
	err = tx.QueryRow("SELECT username FROM users WHERE id = $1", rental.UserID).Scan(&rental.UserUsername)
	if err != nil {
		r.logger.Error("Failed to get user details", zap.Int64("userID", rental.UserID), zap.Error(err))
		return nil, err
	}

	err = tx.QueryRow("SELECT title, author FROM books WHERE id = $1", bookID).Scan(&rental.BookTitle, &rental.BookAuthor)
	if err != nil {
		r.logger.Error("Failed to get book details", zap.Int64("bookID", bookID), zap.Error(err))
		return nil, err
	}

	rental.BookID = bookID

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	return &rental, nil
}

// Extend extends the due date of a rental
func (r *RentalRepository) Extend(id int64, newDueDate time.Time) (*domain.Rental, error) {
	query := `
		UPDATE rentals
		SET due_date = $2, updated_at = NOW()
		WHERE id = $1 AND status = 'active'
		RETURNING id, user_id, book_id, rental_date, due_date, return_date, status, created_at, updated_at
	`

	var rental domain.Rental
	var returnDate sql.NullTime

	err := r.db.QueryRow(query, id, newDueDate).Scan(
		&rental.ID,
		&rental.UserID,
		&rental.BookID,
		&rental.RentalDate,
		&rental.DueDate,
		&returnDate,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Check if rental exists
			exists, _ := r.rentalExists(id)
			if !exists {
				return nil, domain.ErrRentalNotFound
			}
			// Rental exists but not active
			return nil, domain.ErrRentalNotActive
		}
		r.logger.Error("Failed to extend rental", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if returnDate.Valid {
		rental.ReturnDate = &returnDate.Time
	}

	// Get user and book details
	err = r.db.QueryRow("SELECT username FROM users WHERE id = $1", rental.UserID).Scan(&rental.UserUsername)
	if err != nil {
		r.logger.Error("Failed to get user details", zap.Int64("userID", rental.UserID), zap.Error(err))
		return nil, err
	}

	err = r.db.QueryRow("SELECT title, author FROM books WHERE id = $1", rental.BookID).Scan(&rental.BookTitle, &rental.BookAuthor)
	if err != nil {
		r.logger.Error("Failed to get book details", zap.Int64("bookID", rental.BookID), zap.Error(err))
		return nil, err
	}

	return &rental, nil
}

// Delete deletes a rental
func (r *RentalRepository) Delete(id int64) error {
	query := `DELETE FROM rentals WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete rental", zap.Int64("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrRentalNotFound
	}

	return nil
}

// Helper methods

// queryRentals executes a query and returns a list of rentals
func (r *RentalRepository) queryRentals(query string, args ...interface{}) ([]*domain.Rental, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.Error("Failed to query rentals", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rentals []*domain.Rental
	for rows.Next() {
		var rental domain.Rental
		var returnDate sql.NullTime

		err := rows.Scan(
			&rental.ID,
			&rental.UserID,
			&rental.BookID,
			&rental.RentalDate,
			&rental.DueDate,
			&returnDate,
			&rental.Status,
			&rental.CreatedAt,
			&rental.UpdatedAt,
			&rental.UserUsername,
			&rental.BookTitle,
			&rental.BookAuthor,
		)
		if err != nil {
			r.logger.Error("Failed to scan rental row", zap.Error(err))
			return nil, err
		}

		if returnDate.Valid {
			rental.ReturnDate = &returnDate.Time
		}

		rentals = append(rentals, &rental)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating rental rows", zap.Error(err))
		return nil, err
	}

	return rentals, nil
}

// rentalExists checks if a rental exists
func (r *RentalRepository) rentalExists(id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rentals WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check if rental exists", zap.Int64("id", id), zap.Error(err))
		return false, err
	}
	return exists, nil
}
