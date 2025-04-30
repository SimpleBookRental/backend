package service

import (
	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/internal/repository"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
)

// Service is a factory for all services
type Service struct {
	User     domain.UserService
	Auth     AuthService
	Category domain.CategoryService
	Book     domain.BookService
	Rental   domain.RentalService
	Payment  domain.PaymentService
	Report   ReportService
	Logger   *logger.Logger
}

// NewService creates a new service factory
func NewService(repo *repository.Repository, cfg *config.Config, jwtService *auth.JWTService, logger *logger.Logger) *Service {
	serviceLogger := logger.Named("service")

	userService := NewUserService(repo.User, serviceLogger.Named("user"))
	authService := NewAuthService(repo.User, jwtService, serviceLogger.Named("auth"))
	categoryService := NewCategoryService(repo.Category, serviceLogger.Named("category"))
	bookService := NewBookService(repo.Book, repo.Category, serviceLogger.Named("book"))
	rentalService := NewRentalService(repo.Rental, repo.Book, cfg.Rental, serviceLogger.Named("rental"))
	paymentService := NewPaymentService(repo.Payment, repo.Rental, serviceLogger.Named("payment"))
	reportService := NewReportService(repo.Book, repo.Rental, repo.Payment, serviceLogger.Named("report"))

	return &Service{
		User:     userService,
		Auth:     authService,
		Category: categoryService,
		Book:     bookService,
		Rental:   rentalService,
		Payment:  paymentService,
		Report:   reportService,
		Logger:   serviceLogger,
	}
}

// AuthService defines the interface for authentication service
type AuthService interface {
	Register(user *domain.User, password string) (*domain.User, error)
	Login(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(token string) error
}

// ReportService defines the interface for report service
type ReportService interface {
	GetPopularBooks(limit, offset int32) ([]*domain.Book, error)
	GetRevenueReport(startDate, endDate string) ([]*domain.RevenueReport, error)
	GetOverdueBooks(limit, offset int32) ([]*domain.Rental, error)
}
