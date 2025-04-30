package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// PaymentRepository implements domain.PaymentRepository
type PaymentRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewPaymentRepository creates a new PaymentRepository
func NewPaymentRepository(conn *DBConn, logger *logger.Logger) domain.PaymentRepository {
	return &PaymentRepository{
		db:     conn.DB,
		logger: logger,
	}
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(id int64) (*domain.Payment, error) {
	query := `
		SELECT p.id, p.user_id, p.rental_id, p.amount, p.payment_date, p.payment_method, p.status,
			   p.transaction_id, p.created_at, p.updated_at, u.username as user_username,
			   b.title as book_title
		FROM payments p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN rentals r ON p.rental_id = r.id
		LEFT JOIN books b ON r.book_id = b.id
		WHERE p.id = $1
	`

	var payment domain.Payment
	var rentalID sql.NullInt64
	var bookTitle sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&rentalID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.UserUsername,
		&bookTitle,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		r.logger.Error("Failed to get payment by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if rentalID.Valid {
		payment.RentalID = &rentalID.Int64
	}

	if bookTitle.Valid {
		payment.BookTitle = bookTitle.String
	}

	return &payment, nil
}

// List retrieves a list of payments with pagination
func (r *PaymentRepository) List(limit, offset int32) ([]*domain.Payment, error) {
	query := `
		SELECT p.id, p.user_id, p.rental_id, p.amount, p.payment_date, p.payment_method, p.status,
			   p.transaction_id, p.created_at, p.updated_at, u.username as user_username,
			   b.title as book_title
		FROM payments p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN rentals r ON p.rental_id = r.id
		LEFT JOIN books b ON r.book_id = b.id
		ORDER BY p.payment_date DESC
		LIMIT $1 OFFSET $2
	`

	return r.queryPayments(query, limit, offset)
}

// ListByUser retrieves a list of payments for a specific user with pagination
func (r *PaymentRepository) ListByUser(userID int64, limit, offset int32) ([]*domain.Payment, error) {
	query := `
		SELECT p.id, p.user_id, p.rental_id, p.amount, p.payment_date, p.payment_method, p.status,
			   p.transaction_id, p.created_at, p.updated_at, u.username as user_username,
			   b.title as book_title
		FROM payments p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN rentals r ON p.rental_id = r.id
		LEFT JOIN books b ON r.book_id = b.id
		WHERE p.user_id = $1
		ORDER BY p.payment_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list payments by user", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var payment domain.Payment
		var rentalID sql.NullInt64
		var bookTitle sql.NullString

		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&rentalID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.PaymentMethod,
			&payment.Status,
			&payment.TransactionID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.UserUsername,
			&bookTitle,
		)
		if err != nil {
			r.logger.Error("Failed to scan payment row", zap.Error(err))
			return nil, err
		}

		if rentalID.Valid {
			payment.RentalID = &rentalID.Int64
		}

		if bookTitle.Valid {
			payment.BookTitle = bookTitle.String
		}

		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating payment rows", zap.Error(err))
		return nil, err
	}

	return payments, nil
}

// ListByRental retrieves a list of payments for a specific rental
func (r *PaymentRepository) ListByRental(rentalID int64) ([]*domain.Payment, error) {
	query := `
		SELECT p.id, p.user_id, p.rental_id, p.amount, p.payment_date, p.payment_method, p.status,
			   p.transaction_id, p.created_at, p.updated_at, u.username as user_username,
			   b.title as book_title
		FROM payments p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN rentals r ON p.rental_id = r.id
		LEFT JOIN books b ON r.book_id = b.id
		WHERE p.rental_id = $1
		ORDER BY p.payment_date DESC
	`

	rows, err := r.db.Query(query, rentalID)
	if err != nil {
		r.logger.Error("Failed to list payments by rental", zap.Int64("rentalID", rentalID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var payment domain.Payment
		var paymentRentalID sql.NullInt64
		var bookTitle sql.NullString

		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&paymentRentalID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.PaymentMethod,
			&payment.Status,
			&payment.TransactionID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.UserUsername,
			&bookTitle,
		)
		if err != nil {
			r.logger.Error("Failed to scan payment row", zap.Error(err))
			return nil, err
		}

		if paymentRentalID.Valid {
			payment.RentalID = &paymentRentalID.Int64
		}

		if bookTitle.Valid {
			payment.BookTitle = bookTitle.String
		}

		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating payment rows", zap.Error(err))
		return nil, err
	}

	return payments, nil
}

// Create creates a new payment
func (r *PaymentRepository) Create(payment *domain.Payment) (*domain.Payment, error) {
	query := `
		INSERT INTO payments (user_id, rental_id, amount, payment_date, payment_method, status, transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, rental_id, amount, payment_date, payment_method, status, transaction_id, created_at, updated_at
	`

	var rentalID sql.NullInt64
	if payment.RentalID != nil {
		rentalID.Int64 = *payment.RentalID
		rentalID.Valid = true
	}

	err := r.db.QueryRow(
		query,
		payment.UserID,
		rentalID,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
		payment.Status,
		payment.TransactionID,
	).Scan(
		&payment.ID,
		&payment.UserID,
		&rentalID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create payment", zap.Error(err))
		return nil, err
	}

	if rentalID.Valid {
		payment.RentalID = &rentalID.Int64
	}

	// Get user details
	err = r.db.QueryRow("SELECT username FROM users WHERE id = $1", payment.UserID).Scan(&payment.UserUsername)
	if err != nil {
		r.logger.Error("Failed to get user details", zap.Int64("userID", payment.UserID), zap.Error(err))
		return nil, err
	}

	// Get book title if rental is provided
	if payment.RentalID != nil {
		var bookTitle sql.NullString
		err = r.db.QueryRow(`
			SELECT b.title
			FROM rentals r
			JOIN books b ON r.book_id = b.id
			WHERE r.id = $1
		`, *payment.RentalID).Scan(&bookTitle)
		if err == nil && bookTitle.Valid {
			payment.BookTitle = bookTitle.String
		}
	}

	return payment, nil
}

// UpdateStatus updates the status of a payment
func (r *PaymentRepository) UpdateStatus(id int64, status domain.PaymentStatus) (*domain.Payment, error) {
	query := `
		UPDATE payments
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, rental_id, amount, payment_date, payment_method, status, transaction_id, created_at, updated_at
	`

	var payment domain.Payment
	var rentalID sql.NullInt64

	err := r.db.QueryRow(query, id, status).Scan(
		&payment.ID,
		&payment.UserID,
		&rentalID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		r.logger.Error("Failed to update payment status", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if rentalID.Valid {
		payment.RentalID = &rentalID.Int64
	}

	// Get user details
	err = r.db.QueryRow("SELECT username FROM users WHERE id = $1", payment.UserID).Scan(&payment.UserUsername)
	if err != nil {
		r.logger.Error("Failed to get user details", zap.Int64("userID", payment.UserID), zap.Error(err))
		return nil, err
	}

	// Get book title if rental is provided
	if payment.RentalID != nil {
		var bookTitle sql.NullString
		err = r.db.QueryRow(`
			SELECT b.title
			FROM rentals r
			JOIN books b ON r.book_id = b.id
			WHERE r.id = $1
		`, *payment.RentalID).Scan(&bookTitle)
		if err == nil && bookTitle.Valid {
			payment.BookTitle = bookTitle.String
		}
	}

	return &payment, nil
}

// Delete deletes a payment
func (r *PaymentRepository) Delete(id int64) error {
	query := `DELETE FROM payments WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete payment", zap.Int64("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}

// GetRevenueReport generates a revenue report for a specific time period
func (r *PaymentRepository) GetRevenueReport(startDate, endDate time.Time) ([]*domain.RevenueReport, error) {
	query := `
		SELECT 
			DATE_TRUNC('month', payment_date) as month,
			SUM(amount) as total_revenue,
			COUNT(*) as payment_count
		FROM payments
		WHERE status = 'completed'
			AND payment_date BETWEEN $1 AND $2
		GROUP BY DATE_TRUNC('month', payment_date)
		ORDER BY month
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		r.logger.Error("Failed to get revenue report", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var reports []*domain.RevenueReport
	for rows.Next() {
		var report domain.RevenueReport
		err := rows.Scan(
			&report.Month,
			&report.TotalRevenue,
			&report.PaymentCount,
		)
		if err != nil {
			r.logger.Error("Failed to scan revenue report row", zap.Error(err))
			return nil, err
		}
		reports = append(reports, &report)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating revenue report rows", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

// Helper methods

// queryPayments executes a query and returns a list of payments
func (r *PaymentRepository) queryPayments(query string, args ...interface{}) ([]*domain.Payment, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.Error("Failed to query payments", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var payment domain.Payment
		var rentalID sql.NullInt64
		var bookTitle sql.NullString

		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&rentalID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.PaymentMethod,
			&payment.Status,
			&payment.TransactionID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.UserUsername,
			&bookTitle,
		)
		if err != nil {
			r.logger.Error("Failed to scan payment row", zap.Error(err))
			return nil, err
		}

		if rentalID.Valid {
			payment.RentalID = &rentalID.Int64
		}

		if bookTitle.Valid {
			payment.BookTitle = bookTitle.String
		}

		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating payment rows", zap.Error(err))
		return nil, err
	}

	return payments, nil
}
