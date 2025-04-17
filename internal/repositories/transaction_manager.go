package repositories

import (
	"fmt"

	"gorm.io/gorm"
)

// TransactionManagerInterface defines transaction manager behavior for mocking
type TransactionManagerInterface interface {
	WithTransaction(fn func(tx *gorm.DB) error) error
}

// TransactionManager manages database transactions
type TransactionManager struct {
	db *gorm.DB
}

// Ensure TransactionManager implements TransactionManagerInterface
var _ TransactionManagerInterface = (*TransactionManager)(nil)

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// Begin starts a new transaction
func (tm *TransactionManager) Begin() *gorm.DB {
	return tm.db.Begin()
}

// WithTransaction executes a function within a transaction
func (tm *TransactionManager) WithTransaction(fn func(tx *gorm.DB) error) error {
	tx := tm.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
