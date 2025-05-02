package service

import (
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// RentalServiceImpl implements domain.RentalService
type RentalServiceImpl struct {
	repo       domain.RentalRepository
	bookRepo   domain.BookRepository
	config     config.RentalConfig
	logger     *logger.Logger
}

// NewRentalService creates a new RentalService
func NewRentalService(repo domain.RentalRepository, bookRepo domain.BookRepository, config config.RentalConfig, logger *logger.Logger) domain.RentalService {
	return &RentalServiceImpl{
		repo:       repo,
		bookRepo:   bookRepo,
		config:     config,
		logger:     logger,
	}
}

// GetByID retrieves a rental by ID
func (s *RentalServiceImpl) GetByID(id int64) (*domain.Rental, error) {
	rental, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// Update status if it's overdue but not marked as overdue
	if rental.Status == domain.RentalStatusActive && s.IsOverdue(rental) {
		rental.Status = domain.RentalStatusOverdue
		rental, err = s.repo.UpdateStatus(rental.ID, domain.RentalStatusOverdue)
		if err != nil {
			s.logger.Error("Failed to update rental status to overdue", zap.Int64("id", id), zap.Error(err))
			// Continue anyway, we'll return the rental with the updated status even if the DB update failed
		}
	}

	return rental, nil
}

// List retrieves a list of rentals with pagination
func (s *RentalServiceImpl) List(limit, offset int32) ([]*domain.Rental, error) {
	rentals, err := s.repo.List(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list rentals", zap.Error(err))
		return nil, err
	}

	// Update status for any overdue rentals
	s.updateOverdueStatus(rentals)

	return rentals, nil
}

// ListByUser retrieves a list of rentals for a specific user with pagination
func (s *RentalServiceImpl) ListByUser(userID int64, limit, offset int32) ([]*domain.Rental, error) {
	rentals, err := s.repo.ListByUser(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list rentals by user", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}

	// Update status for any overdue rentals
	s.updateOverdueStatus(rentals)

	return rentals, nil
}

// ListByBook retrieves a list of rentals for a specific book with pagination
func (s *RentalServiceImpl) ListByBook(bookID int64, limit, offset int32) ([]*domain.Rental, error) {
	rentals, err := s.repo.ListByBook(bookID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list rentals by book", zap.Int64("bookID", bookID), zap.Error(err))
		return nil, err
	}

	// Update status for any overdue rentals
	s.updateOverdueStatus(rentals)

	return rentals, nil
}

// ListActive retrieves a list of active rentals with pagination
func (s *RentalServiceImpl) ListActive(limit, offset int32) ([]*domain.Rental, error) {
	rentals, err := s.repo.ListActive(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list active rentals", zap.Error(err))
		return nil, err
	}

	// Update status for any overdue rentals
	s.updateOverdueStatus(rentals)

	return rentals, nil
}

// ListOverdue retrieves a list of overdue rentals with pagination
func (s *RentalServiceImpl) ListOverdue(limit, offset int32) ([]*domain.Rental, error) {
	// Get current overdue rentals
	rentals, err := s.repo.ListOverdue(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list overdue rentals", zap.Error(err))
		return nil, err
	}

	// Check active rentals to see if any have become overdue
	activeRentals, err := s.repo.ListActive(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list active rentals to check for overdue", zap.Error(err))
		// Continue with the overdue rentals we already have
	} else {
		// Find any active rentals that are now overdue
		var newOverdueRentals []*domain.Rental
		for _, rental := range activeRentals {
			if rental.DueDate.Before(time.Now()) {
				rental.Status = domain.RentalStatusOverdue
				updatedRental, err := s.repo.UpdateStatus(rental.ID, domain.RentalStatusOverdue)
				if err != nil {
					s.logger.Error("Failed to update rental status to overdue", 
						zap.Int64("id", rental.ID), zap.Error(err))
					// Continue with the next rental
					continue
				}
				newOverdueRentals = append(newOverdueRentals, updatedRental)
			}
		}

		// Add any newly found overdue rentals to the list
		if len(newOverdueRentals) > 0 {
			rentals = append(rentals, newOverdueRentals...)
			// We might exceed the limit now, so trim the list if needed
			if int32(len(rentals)) > limit {
				rentals = rentals[:limit]
			}
		}
	}

	return rentals, nil
}

// Create creates a new rental
func (s *RentalServiceImpl) Create(rental *domain.Rental) (*domain.Rental, error) {
	// Check if book exists and is available
	isAvailable, err := s.IsBookAvailable(rental.BookID)
	if err != nil {
		s.logger.Error("Failed to check book availability", zap.Int64("bookID", rental.BookID), zap.Error(err))
		return nil, err
	}

	if !isAvailable {
		return nil, domain.ErrBookNotAvailable
	}

	// Set rental date to now if not provided
	if rental.RentalDate.IsZero() {
		rental.RentalDate = time.Now()
	}

	// Set due date based on rental date and default rental days if not provided
	if rental.DueDate.IsZero() {
		rental.DueDate = rental.RentalDate.AddDate(0, 0, s.config.DefaultRentalDays)
	}

	// Set status to active if not provided
	if rental.Status == "" {
		rental.Status = domain.RentalStatusActive
	}

	// Create rental
	createdRental, err := s.repo.Create(rental)
	if err != nil {
		s.logger.Error("Failed to create rental", zap.Error(err))
		return nil, err
	}

	return createdRental, nil
}

// Return processes the return of a rental
func (s *RentalServiceImpl) Return(id int64) (*domain.Rental, error) {
	// Check if rental exists and is active or overdue
	rental, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if rental.Status != domain.RentalStatusActive && rental.Status != domain.RentalStatusOverdue {
		return nil, domain.ErrRentalNotActive
	}

	// Return rental
	returnedRental, err := s.repo.Return(id)
	if err != nil {
		s.logger.Error("Failed to return rental", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return returnedRental, nil
}

// Extend extends the due date of a rental
func (s *RentalServiceImpl) Extend(id int64, days int) (*domain.Rental, error) {
	// Check if rental exists and is active
	rental, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if rental.Status != domain.RentalStatusActive {
		return nil, domain.ErrRentalNotActive
	}

	// Validate extension days
	if days <= 0 {
		return nil, domain.NewInvalidInputError("extension days must be positive")
	}

	if days > s.config.MaxRentalExtensionDays {
		return nil, domain.NewInvalidInputError("extension days cannot exceed maximum allowed")
	}

	// Calculate new due date
	newDueDate := rental.DueDate.AddDate(0, 0, days)

	// Extend rental
	extendedRental, err := s.repo.Extend(id, newDueDate)
	if err != nil {
		s.logger.Error("Failed to extend rental", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return extendedRental, nil
}

// CalculateLateFee calculates the late fee for a rental
func (s *RentalServiceImpl) CalculateLateFee(rental *domain.Rental) (float64, error) {
	// If rental is not overdue, no late fee
	if !s.IsOverdue(rental) {
		return 0, nil
	}

	// If rental is returned, calculate late fee based on return date
	var lateDays int
	if rental.ReturnDate != nil {
		lateDays = int(rental.ReturnDate.Sub(rental.DueDate).Hours() / 24)
	} else {
		// If rental is not returned, calculate late fee based on current date
		lateDays = int(time.Now().Sub(rental.DueDate).Hours() / 24)
	}

	// Ensure late days is at least 1
	if lateDays < 1 {
		lateDays = 1
	}

	return float64(lateDays) * s.config.LateFeePerDay, nil
}

// IsOverdue checks if a rental is overdue
func (s *RentalServiceImpl) IsOverdue(rental *domain.Rental) bool {
	// If rental is not active, it's not overdue
	if rental.Status != domain.RentalStatusActive {
		return false
	}

	// If due date is in the future, it's not overdue
	return rental.DueDate.Before(time.Now())
}

// Helper methods

// IsBookAvailable checks if a book is available for rental
func (s *RentalServiceImpl) IsBookAvailable(bookID int64) (bool, error) {
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return false, err
	}
	return book.AvailableCopies > 0, nil
}

// updateOverdueStatus updates the status of any overdue rentals in the given list
func (s *RentalServiceImpl) updateOverdueStatus(rentals []*domain.Rental) {
	for _, rental := range rentals {
		if rental.Status == domain.RentalStatusActive && s.IsOverdue(rental) {
			rental.Status = domain.RentalStatusOverdue
			_, err := s.repo.UpdateStatus(rental.ID, domain.RentalStatusOverdue)
			if err != nil {
				s.logger.Error("Failed to update rental status to overdue", 
					zap.Int64("id", rental.ID), zap.Error(err))
				// Continue with the next rental
			}
		}
	}
}
