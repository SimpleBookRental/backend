package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupTxMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestTransactionManager_WithTransaction_Success(t *testing.T) {
	gdb, mock, cleanup := setupTxMockDB(t)
	defer cleanup()
	tm := NewTransactionManager(gdb)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := tm.WithTransaction(func(tx *gorm.DB) error { return nil })

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionManager_WithTransaction_ErrorRollback(t *testing.T) {
	gdb, mock, cleanup := setupTxMockDB(t)
	defer cleanup()
	tm := NewTransactionManager(gdb)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := tm.WithTransaction(func(tx *gorm.DB) error {
		return errors.New("fail in tx")
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionManager_WithTransaction_PanicRollback(t *testing.T) {
	gdb, mock, cleanup := setupTxMockDB(t)
	defer cleanup()
	tm := NewTransactionManager(gdb)

	mock.ExpectBegin()
	mock.ExpectRollback()

	assert.Panics(t, func() {
		_ = tm.WithTransaction(func(tx *gorm.DB) error {
			panic("panic in tx")
		})
	})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionManager_WithTransaction_CommitError(t *testing.T) {
	gdb, mock, cleanup := setupTxMockDB(t)
	defer cleanup()
	tm := NewTransactionManager(gdb)

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	// You DON'T need mock.ExpectRollback() anymore, 
	// because the new defer calls rollback only if err is set. 
	// In this test, err was already returned early → Rollback() won't be called again.

	err := tm.WithTransaction(func(tx *gorm.DB) error { return nil })

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
