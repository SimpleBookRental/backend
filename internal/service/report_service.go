package service

import (
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// ReportServiceImpl implements ReportService
type ReportServiceImpl struct {
	bookRepo    domain.BookRepository
	rentalRepo  domain.RentalRepository
	paymentRepo domain.PaymentRepository
	logger      *logger.Logger
}

// NewReportService creates a new ReportService
func NewReportService(bookRepo domain.BookRepository, rentalRepo domain.RentalRepository, paymentRepo domain.PaymentRepository, logger *logger.Logger) ReportService {
	return &ReportServiceImpl{
		bookRepo:    bookRepo,
		rentalRepo:  rentalRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

// GetPopularBooks retrieves a list of popular books based on rental frequency
func (s *ReportServiceImpl) GetPopularBooks(limit, offset int32) ([]*domain.Book, error) {
	// In a real implementation, this would use a more efficient query
	// that joins books and rentals tables and counts rentals per book
	// For this implementation, we'll use a simplistic approach
	
	// Get all rentals
	rentals, err := s.rentalRepo.List(1000, 0) // Using a larger limit to get sufficient data
	if err != nil {
		s.logger.Error("Failed to list rentals for popular books report", zap.Error(err))
		return nil, err
	}
	
	// Count rental frequency for each book
	bookRentalCount := make(map[int64]int)
	for _, rental := range rentals {
		bookRentalCount[rental.BookID]++
	}
	
	// Sort books by rental count (using a simple approach for demonstration)
	type bookRank struct {
		bookID int64
		count  int
	}
	
	var rankedBooks []bookRank
	for bookID, count := range bookRentalCount {
		rankedBooks = append(rankedBooks, bookRank{bookID, count})
	}
	
	// Simple bubble sort by count (descending)
	for i := 0; i < len(rankedBooks)-1; i++ {
		for j := 0; j < len(rankedBooks)-i-1; j++ {
			if rankedBooks[j].count < rankedBooks[j+1].count {
				rankedBooks[j], rankedBooks[j+1] = rankedBooks[j+1], rankedBooks[j]
			}
		}
	}
	
	// Apply pagination
	start := int(offset)
	end := int(offset + limit)
	if start >= len(rankedBooks) {
		return []*domain.Book{}, nil
	}
	if end > len(rankedBooks) {
		end = len(rankedBooks)
	}
	
	pagedRankBooks := rankedBooks[start:end]
	
	// Get book details for the popular books
	var popularBooks []*domain.Book
	for _, rb := range pagedRankBooks {
		book, err := s.bookRepo.GetByID(rb.bookID)
		if err != nil {
			s.logger.Error("Failed to get book details for popular book", 
				zap.Int64("bookID", rb.bookID), zap.Error(err))
			continue
		}
		popularBooks = append(popularBooks, book)
	}
	
	return popularBooks, nil
}

// GetRevenueReport retrieves a revenue report for a specific time period
func (s *ReportServiceImpl) GetRevenueReport(startDate, endDate string) ([]*domain.RevenueReport, error) {
	// Parse date strings to time.Time
	start, err := parseDate(startDate)
	if err != nil {
		s.logger.Error("Failed to parse start date", zap.String("startDate", startDate), zap.Error(err))
		return nil, domain.NewInvalidInputError("invalid start date format, use YYYY-MM-DD")
	}
	
	end, err := parseDate(endDate)
	if err != nil {
		s.logger.Error("Failed to parse end date", zap.String("endDate", endDate), zap.Error(err))
		return nil, domain.NewInvalidInputError("invalid end date format, use YYYY-MM-DD")
	}
	
	// Add 1 day to end date to include the full day
	end = end.AddDate(0, 0, 1)
	
	// Get revenue report from payment repository
	report, err := s.paymentRepo.GetRevenueReport(start, end)
	if err != nil {
		s.logger.Error("Failed to get revenue report", 
			zap.Time("startDate", start), 
			zap.Time("endDate", end), 
			zap.Error(err))
		return nil, err
	}
	
	return report, nil
}

// GetOverdueBooks retrieves a list of overdue rentals
func (s *ReportServiceImpl) GetOverdueBooks(limit, offset int32) ([]*domain.Rental, error) {
	// This leverages the existing overdue rentals functionality
	rentals, err := s.rentalRepo.ListOverdue(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list overdue rentals", zap.Error(err))
		return nil, err
	}
	
	// Check active rentals to find any additional overdue rentals
	// that haven't been updated in the database yet
	activeRentals, err := s.rentalRepo.ListActive(100, 0)
	if err != nil {
		s.logger.Error("Failed to list active rentals", zap.Error(err))
		// Continue with just the known overdue rentals
	} else {
		now := time.Now()
		var newOverdueRentals []*domain.Rental
		
		for _, rental := range activeRentals {
			if rental.DueDate.Before(now) {
				// Mark as overdue in the database
				rental.Status = domain.RentalStatusOverdue
				_, updateErr := s.rentalRepo.UpdateStatus(rental.ID, domain.RentalStatusOverdue)
				if updateErr != nil {
					s.logger.Error("Failed to update rental status", 
						zap.Int64("id", rental.ID), zap.Error(updateErr))
					// Continue anyway, we want to include this rental in the report
				}
				newOverdueRentals = append(newOverdueRentals, rental)
			}
		}
		
		// Add newly discovered overdue rentals to the result
		if len(newOverdueRentals) > 0 {
			rentals = append(rentals, newOverdueRentals...)
			// Apply pagination if needed
			if int32(len(rentals)) > limit {
				rentals = rentals[:limit]
			}
		}
	}
	
	return rentals, nil
}

// Helper functions

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		// Default to beginning of current month if not provided
		now := time.Now()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), nil
	}
	
	// Parse date in YYYY-MM-DD format
	return time.Parse("2006-01-02", dateStr)
}
