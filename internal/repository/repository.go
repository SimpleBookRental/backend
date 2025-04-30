package repository

import (
	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
)

// Repository is a factory for all repositories
type Repository struct {
	User     domain.UserRepository
	Category domain.CategoryRepository
	Book     domain.BookRepository
	Rental   domain.RentalRepository
	Payment  domain.PaymentRepository
	Logger   *logger.Logger
}

// NewRepository creates a new repository factory
func NewRepository(conn *DBConn) *Repository {
	logger := conn.Logger.Named("repository")

	return &Repository{
		User:     NewUserRepository(conn, logger.Named("user")),
		Category: NewCategoryRepository(conn, logger.Named("category")),
		Book:     NewBookRepository(conn, logger.Named("book")),
		Rental:   NewRentalRepository(conn, logger.Named("rental")),
		Payment:  NewPaymentRepository(conn, logger.Named("payment")),
		Logger:   logger,
	}
}
