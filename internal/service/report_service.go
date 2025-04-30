package service

import (
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// ReportServiceImpl implements ReportService
type ReportServiceImpl struct {
	bookRepo   domain.BookRepository
	rentalRepo domain.RentalRepository
	paymentRepo domain.PaymentRepository
	logger     *logger.Logger
}

// NewReportService creates a new ReportService
func NewReportService(bookRepo domain.BookRepository, rentalRepo domain.RentalRepository, paymentRepo domain.PaymentRepository, logger *logger.Logger) ReportService {
	return &ReportServiceImpl{
		bookRepo:   bookRepo,
		rentalRepo: rentalRepo,
		paymentRepo: paymentRepo,
		logger:     logger,
	}
}

// GetPopularBooks retrieves a list of popular books based on rental count
func (s *ReportServiceImpl) GetPopularBooks(limit, offset int32) ([]*domain.Book, error) {
	// In a real-world application, this would be a more complex query
	// For simplicity, we'll just return a list of books
	books, err := s.bookRepo.List(limit, offset)
	if err != nil {
		s.logger.Error("Failed to get popular books", zap.Error(err))
		return nil, err
	}
	return books, nil
}

// GetRevenueReport generates a revenue report for a specific time period
func (s *ReportServiceImpl) GetRevenueReport(startDate, endDate string) ([]*domain.RevenueReport, error) {
	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.logger.Error("Failed to parse start date", zap.String("startDate", startDate), zap.Error(err))
		return nil, domain.NewInvalidInputError("invalid start date format, expected YYYY-MM-DD")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.logger.Error("Failed to parse end date", zap.String("endDate", endDate), zap.Error(err))
		return nil, domain.NewInvalidInputError("invalid end date format, expected YYYY-MM-DD")
	}

	// Ensure end date is at the end of the day
	end = end.Add(24*time.Hour - time.Second)

	// Get revenue report
	report, err := s.paymentRepo.GetRevenueReport(start, end)
	if err != nil {
		s.logger.Error("Failed to get revenue report", zap.Error(err))
		return nil, err
	}

	return report, nil
}

// GetOverdueBooks retrieves a list of overdue rentals
func (s *ReportServiceImpl) GetOverdueBooks(limit, offset int32) ([]*domain.Rental, error) {
	rentals, err := s.rentalRepo.ListOverdue(limit, offset)
	if err != nil {
		s.logger.Error("Failed to get overdue books", zap.Error(err))
		return nil, err
	}
	return rentals, nil
}
