package repositories

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of *gorm.DB
type MockDB struct {
	mock.Mock
}

// Create mocks the Create method of *gorm.DB
func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// First mocks the First method of *gorm.DB
func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

// Find mocks the Find method of *gorm.DB
func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

// Where mocks the Where method of *gorm.DB
func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

// Save mocks the Save method of *gorm.DB
func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// Delete mocks the Delete method of *gorm.DB
func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(value, conds)
	return args.Get(0).(*gorm.DB)
}

// Model mocks the Model method of *gorm.DB
func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// Updates mocks the Updates method of *gorm.DB
func (m *MockDB) Updates(values interface{}) *gorm.DB {
	args := m.Called(values)
	return args.Get(0).(*gorm.DB)
}

// Error returns the error from the mock DB
func (m *MockDB) Error() error {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

// RowsAffected returns the number of rows affected
func (m *MockDB) RowsAffected() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

// MockGormDB is a mock implementation of gorm.DB
type MockGormDB struct {
	mock.Mock
}

// DB returns a mock DB
func (m *MockGormDB) DB() *MockDB {
	args := m.Called()
	return args.Get(0).(*MockDB)
}

// Begin mocks the Begin method of gorm.DB
func (m *MockGormDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Commit mocks the Commit method of gorm.DB
func (m *MockGormDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Rollback mocks the Rollback method of gorm.DB
func (m *MockGormDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
