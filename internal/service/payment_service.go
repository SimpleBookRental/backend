package service

import (
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// PaymentServiceImpl implements domain.PaymentService
type PaymentServiceImpl struct {
	repo      domain.PaymentRepository
	rentalRepo domain.RentalRepository
	logger    *logger.Logger
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(repo domain.PaymentRepository, rentalRepo domain.RentalRepository, logger *logger.Logger) domain.PaymentService {
	return &PaymentServiceImpl{
		repo:      repo,
		rentalRepo: rentalRepo,
		logger:    logger,
	}
}

// GetByID retrieves a payment by ID
func (s *PaymentServiceImpl) GetByID(id int64) (*domain.Payment, error) {
	payment, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get payment by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return payment, nil
}

// List retrieves a list of payments with pagination
func (s *PaymentServiceImpl) List(limit, offset int32) ([]*domain.Payment, error) {
	payments, err := s.repo.List(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list payments", zap.Error(err))
		return nil, err
	}
	return payments, nil
}

// ListByUser retrieves a list of payments for a specific user with pagination
func (s *PaymentServiceImpl) ListByUser(userID int64, limit, offset int32) ([]*domain.Payment, error) {
	payments, err := s.repo.ListByUser(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list payments by user", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}
	return payments, nil
}

// ListByRental retrieves a list of payments for a specific rental
func (s *PaymentServiceImpl) ListByRental(rentalID int64) ([]*domain.Payment, error) {
	payments, err := s.repo.ListByRental(rentalID)
	if err != nil {
		s.logger.Error("Failed to list payments by rental", zap.Int64("rentalID", rentalID), zap.Error(err))
		return nil, err
	}
	return payments, nil
}

// Create creates a new payment
func (s *PaymentServiceImpl) Create(payment *domain.Payment) (*domain.Payment, error) {
	// Validate rental if provided
	if payment.RentalID != nil {
		rental, err := s.rentalRepo.GetByID(*payment.RentalID)
		if err != nil {
			s.logger.Error("Failed to get rental by ID", zap.Int64("rentalID", *payment.RentalID), zap.Error(err))
			return nil, err
		}

		// Ensure payment is associated with the correct user
		if payment.UserID != rental.UserID {
			return nil, domain.NewInvalidInputError("payment user ID does not match rental user ID")
		}
	}

	// Set payment date to now if not provided
	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}

	// Create payment
	createdPayment, err := s.repo.Create(payment)
	if err != nil {
		s.logger.Error("Failed to create payment", zap.Error(err))
		return nil, err
	}

	return createdPayment, nil
}

// ProcessPayment processes a payment
func (s *PaymentServiceImpl) ProcessPayment(payment *domain.Payment) (*domain.Payment, error) {
	// In a real-world application, this would integrate with a payment gateway
	// For simplicity, we'll just create the payment with a completed status

	// Set payment date to now if not provided
	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}

	// Set status to completed
	payment.Status = domain.PaymentStatusCompleted

	// Generate a transaction ID if not provided
	if payment.TransactionID == "" {
		payment.TransactionID = generateTransactionID()
	}

	// Create payment
	processedPayment, err := s.repo.Create(payment)
	if err != nil {
		s.logger.Error("Failed to process payment", zap.Error(err))
		return nil, err
	}

	return processedPayment, nil
}

// RefundPayment refunds a payment
func (s *PaymentServiceImpl) RefundPayment(id int64) (*domain.Payment, error) {
	// Get payment
	payment, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get payment by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// Check if payment can be refunded
	if payment.Status != domain.PaymentStatusCompleted {
		return nil, domain.NewInvalidInputError("only completed payments can be refunded")
	}

	// Update payment status to refunded
	refundedPayment, err := s.repo.UpdateStatus(id, domain.PaymentStatusRefunded)
	if err != nil {
		s.logger.Error("Failed to refund payment", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return refundedPayment, nil
}

// GetRevenueReport generates a revenue report for a specific time period
func (s *PaymentServiceImpl) GetRevenueReport(startDate, endDate time.Time) ([]*domain.RevenueReport, error) {
	// Get revenue report
	report, err := s.repo.GetRevenueReport(startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get revenue report", zap.Error(err))
		return nil, err
	}

	return report, nil
}

// Helper functions

// generateTransactionID generates a unique transaction ID
func generateTransactionID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of the specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
	}
	return string(result)
}
