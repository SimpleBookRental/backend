package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// BookRepository implements domain.BookRepository
type BookRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewBookRepository creates a new BookRepository
func NewBookRepository(conn *DBConn, logger *logger.Logger) domain.BookRepository {
	return &BookRepository{
		db:     conn.DB,
		logger: logger,
	}
}

// GetByID retrieves a book by ID
func (r *BookRepository) GetByID(id int64) (*domain.Book, error) {
	query := `
		SELECT b.id, b.title, b.author, b.isbn, b.description, b.published_year, b.publisher,
			   b.total_copies, b.available_copies, b.category_id, c.name as category_name,
			   b.created_at, b.updated_at
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE b.id = $1
	`

	var book domain.Book
	var categoryID sql.NullInt64
	var categoryName sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&categoryName,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		r.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
	}

	if categoryName.Valid {
		book.CategoryName = categoryName.String
	}

	return &book, nil
}

// GetByISBN retrieves a book by ISBN
func (r *BookRepository) GetByISBN(isbn string) (*domain.Book, error) {
	query := `
		SELECT b.id, b.title, b.author, b.isbn, b.description, b.published_year, b.publisher,
			   b.total_copies, b.available_copies, b.category_id, c.name as category_name,
			   b.created_at, b.updated_at
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE b.isbn = $1
	`

	var book domain.Book
	var categoryID sql.NullInt64
	var categoryName sql.NullString

	err := r.db.QueryRow(query, isbn).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&categoryName,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		r.logger.Error("Failed to get book by ISBN", zap.String("isbn", isbn), zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
	}

	if categoryName.Valid {
		book.CategoryName = categoryName.String
	}

	return &book, nil
}

// List retrieves a list of books with pagination
func (r *BookRepository) List(limit, offset int32) ([]*domain.Book, error) {
	query := `
		SELECT b.id, b.title, b.author, b.isbn, b.description, b.published_year, b.publisher,
			   b.total_copies, b.available_copies, b.category_id, c.name as category_name,
			   b.created_at, b.updated_at
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		ORDER BY b.title
		LIMIT $1 OFFSET $2
	`

	return r.queryBooks(query, limit, offset)
}

// ListByCategory retrieves a list of books by category with pagination
func (r *BookRepository) ListByCategory(categoryID int64, limit, offset int32) ([]*domain.Book, error) {
	query := `
		SELECT b.id, b.title, b.author, b.isbn, b.description, b.published_year, b.publisher,
			   b.total_copies, b.available_copies, b.category_id, c.name as category_name,
			   b.created_at, b.updated_at
		FROM books b
		JOIN categories c ON b.category_id = c.id
		WHERE b.category_id = $1
		ORDER BY b.title
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, categoryID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list books by category", zap.Int64("categoryID", categoryID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		var book domain.Book
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Description,
			&book.PublishedYear,
			&book.Publisher,
			&book.TotalCopies,
			&book.AvailableCopies,
			&categoryID,
			&categoryName,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan book row", zap.Error(err))
			return nil, err
		}

		if categoryID.Valid {
			book.CategoryID = categoryID.Int64
		}

		if categoryName.Valid {
			book.CategoryName = categoryName.String
		}

		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating book rows", zap.Error(err))
		return nil, err
	}

	return books, nil
}

// Search searches for books based on search parameters
func (r *BookRepository) Search(params domain.BookSearchParams) ([]*domain.Book, error) {
	query := `
		SELECT b.id, b.title, b.author, b.isbn, b.description, b.published_year, b.publisher,
			   b.total_copies, b.available_copies, b.category_id, c.name as category_name,
			   b.created_at, b.updated_at
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE 1=1
	`

	var conditions []string
	var args []interface{}
	var argIndex int = 1

	if params.Title != "" {
		conditions = append(conditions, fmt.Sprintf("b.title ILIKE $%d", argIndex))
		args = append(args, "%"+params.Title+"%")
		argIndex++
	}

	if params.Author != "" {
		conditions = append(conditions, fmt.Sprintf("b.author ILIKE $%d", argIndex))
		args = append(args, "%"+params.Author+"%")
		argIndex++
	}

	if params.ISBN != "" {
		conditions = append(conditions, fmt.Sprintf("b.isbn = $%d", argIndex))
		args = append(args, params.ISBN)
		argIndex++
	}

	if params.PublishedYear != 0 {
		conditions = append(conditions, fmt.Sprintf("b.published_year = $%d", argIndex))
		args = append(args, params.PublishedYear)
		argIndex++
	}

	if params.CategoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("b.category_id = $%d", argIndex))
		args = append(args, params.CategoryID)
		argIndex++
	}

	if params.Available {
		conditions = append(conditions, "b.available_copies > 0")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY b.title"

	if params.Limit == 0 {
		params.Limit = 10
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.Error("Failed to search books", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		var book domain.Book
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Description,
			&book.PublishedYear,
			&book.Publisher,
			&book.TotalCopies,
			&book.AvailableCopies,
			&categoryID,
			&categoryName,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan book row", zap.Error(err))
			return nil, err
		}

		if categoryID.Valid {
			book.CategoryID = categoryID.Int64
		}

		if categoryName.Valid {
			book.CategoryName = categoryName.String
		}

		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating book rows", zap.Error(err))
		return nil, err
	}

	return books, nil
}

// Create creates a new book
func (r *BookRepository) Create(book *domain.Book) (*domain.Book, error) {
	query := `
		INSERT INTO books (title, author, isbn, description, published_year, publisher, total_copies, available_copies, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, author, isbn, description, published_year, publisher, total_copies, available_copies, category_id, created_at, updated_at
	`

	var categoryID sql.NullInt64
	if book.CategoryID != 0 {
		categoryID.Int64 = book.CategoryID
		categoryID.Valid = true
	}

	err := r.db.QueryRow(
		query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Description,
		book.PublishedYear,
		book.Publisher,
		book.TotalCopies,
		book.AvailableCopies,
		categoryID,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create book", zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
		// Get category name
		var categoryName string
		err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID.Int64).Scan(&categoryName)
		if err == nil {
			book.CategoryName = categoryName
		}
	}

	return book, nil
}

// Update updates an existing book
func (r *BookRepository) Update(book *domain.Book) (*domain.Book, error) {
	query := `
		UPDATE books
		SET title = $2, author = $3, isbn = $4, description = $5, published_year = $6, 
			publisher = $7, category_id = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, author, isbn, description, published_year, publisher, total_copies, available_copies, category_id, created_at, updated_at
	`

	var categoryID sql.NullInt64
	if book.CategoryID != 0 {
		categoryID.Int64 = book.CategoryID
		categoryID.Valid = true
	}

	err := r.db.QueryRow(
		query,
		book.ID,
		book.Title,
		book.Author,
		book.ISBN,
		book.Description,
		book.PublishedYear,
		book.Publisher,
		categoryID,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		r.logger.Error("Failed to update book", zap.Int64("id", book.ID), zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
		// Get category name
		var categoryName string
		err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID.Int64).Scan(&categoryName)
		if err == nil {
			book.CategoryName = categoryName
		}
	}

	return book, nil
}

// UpdateCopies updates the total and available copies of a book
func (r *BookRepository) UpdateCopies(id int64, totalCopies, availableCopies int32) (*domain.Book, error) {
	query := `
		UPDATE books
		SET total_copies = $2, available_copies = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, author, isbn, description, published_year, publisher, total_copies, available_copies, category_id, created_at, updated_at
	`

	var book domain.Book
	var categoryID sql.NullInt64

	err := r.db.QueryRow(
		query,
		id,
		totalCopies,
		availableCopies,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		r.logger.Error("Failed to update book copies", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
		// Get category name
		var categoryName string
		err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID.Int64).Scan(&categoryName)
		if err == nil {
			book.CategoryName = categoryName
		}
	}

	return &book, nil
}

// DecrementAvailableCopies decrements the available copies of a book
func (r *BookRepository) DecrementAvailableCopies(id int64) (*domain.Book, error) {
	query := `
		UPDATE books
		SET available_copies = available_copies - 1, updated_at = NOW()
		WHERE id = $1 AND available_copies > 0
		RETURNING id, title, author, isbn, description, published_year, publisher, total_copies, available_copies, category_id, created_at, updated_at
	`

	var book domain.Book
	var categoryID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Check if the book exists
			exists, _ := r.bookExists(id)
			if !exists {
				return nil, domain.ErrBookNotFound
			}
			// Book exists but no available copies
			return nil, domain.ErrBookNotAvailable
		}
		r.logger.Error("Failed to decrement available copies", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
		// Get category name
		var categoryName string
		err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID.Int64).Scan(&categoryName)
		if err == nil {
			book.CategoryName = categoryName
		}
	}

	return &book, nil
}

// IncrementAvailableCopies increments the available copies of a book
func (r *BookRepository) IncrementAvailableCopies(id int64) (*domain.Book, error) {
	query := `
		UPDATE books
		SET available_copies = available_copies + 1, updated_at = NOW()
		WHERE id = $1 AND available_copies < total_copies
		RETURNING id, title, author, isbn, description, published_year, publisher, total_copies, available_copies, category_id, created_at, updated_at
	`

	var book domain.Book
	var categoryID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedYear,
		&book.Publisher,
		&book.TotalCopies,
		&book.AvailableCopies,
		&categoryID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Check if the book exists
			exists, _ := r.bookExists(id)
			if !exists {
				return nil, domain.ErrBookNotFound
			}
			// Book exists but all copies are available
			return nil, fmt.Errorf("all copies are already available")
		}
		r.logger.Error("Failed to increment available copies", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if categoryID.Valid {
		book.CategoryID = categoryID.Int64
		// Get category name
		var categoryName string
		err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID.Int64).Scan(&categoryName)
		if err == nil {
			book.CategoryName = categoryName
		}
	}

	return &book, nil
}

// Delete deletes a book
func (r *BookRepository) Delete(id int64) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete book", zap.Int64("id", id), zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrBookNotFound
	}

	return nil
}

// Helper methods

// queryBooks executes a query and returns a list of books
func (r *BookRepository) queryBooks(query string, args ...interface{}) ([]*domain.Book, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.Error("Failed to query books", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		var book domain.Book
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Description,
			&book.PublishedYear,
			&book.Publisher,
			&book.TotalCopies,
			&book.AvailableCopies,
			&categoryID,
			&categoryName,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan book row", zap.Error(err))
			return nil, err
		}

		if categoryID.Valid {
			book.CategoryID = categoryID.Int64
		}

		if categoryName.Valid {
			book.CategoryName = categoryName.String
		}

		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating book rows", zap.Error(err))
		return nil, err
	}

	return books, nil
}

// bookExists checks if a book exists
func (r *BookRepository) bookExists(id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check if book exists", zap.Int64("id", id), zap.Error(err))
		return false, err
	}
	return exists, nil
}
