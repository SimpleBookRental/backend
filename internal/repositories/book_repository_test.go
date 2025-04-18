package repositories

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupBookMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestBookRepository_FindByISBN_Found(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "description", "user_id"}).
		AddRow("b1a2c3d4-e5f6-7890-abcd-1234567890ab", "The Great Gatsby", "F. Scott Fitzgerald", "978-1234567890", "desc", "11111111-1111-1111-1111-111111111111")

	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("978-1234567890", sqlmock.AnyArg()).
		WillReturnRows(rows)

	book, err := repo.FindByISBN("978-1234567890")
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "The Great Gatsby", book.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_FindByISBN_NotFound(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("978-1234567890", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "isbn", "description", "user_id"}))

	book, err := repo.FindByISBN("978-1234567890")
	assert.NoError(t, err)
	assert.Nil(t, book)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_FindByISBN_Error(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("978-1234567890", sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	book, err := repo.FindByISBN("978-1234567890")
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_FindByID_Found(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "description", "user_id"}).
		AddRow("b1a2c3d4-e5f6-7890-abcd-1234567890ab", "The Great Gatsby", "F. Scott Fitzgerald", "978-1234567890", "desc", "11111111-1111-1111-1111-111111111111")

	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("b1a2c3d4-e5f6-7890-abcd-1234567890ab", sqlmock.AnyArg()).
		WillReturnRows(rows)
	// Mock preload User
	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
			AddRow("11111111-1111-1111-1111-111111111111", "User1", "user1@example.com", "hashed", "USER"))

	book, err := repo.FindByID("b1a2c3d4-e5f6-7890-abcd-1234567890ab")
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "The Great Gatsby", book.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_FindByID_NotFound(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("b1a2c3d4-e5f6-7890-abcd-1234567890ab", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "isbn", "description", "user_id"}))

	book, err := repo.FindByID("b1a2c3d4-e5f6-7890-abcd-1234567890ab")
	assert.NoError(t, err)
	assert.Nil(t, book)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_FindByID_Error(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("b1a2c3d4-e5f6-7890-abcd-1234567890ab", sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	book, err := repo.FindByID("b1a2c3d4-e5f6-7890-abcd-1234567890ab")
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_Create(t *testing.T) {
	gdb, mock, cleanup := setupBookMockDB(t)
	defer cleanup()
	repo := NewBookRepository(gdb)

	book := &models.Book{
		ID:     "b1a2c3d4-e5f6-7890-abcd-1234567890ab",
		Title:  "The Great Gatsby",
		Author: "F. Scott Fitzgerald",
		ISBN:   "978-1234567890",
		UserID: "11111111-1111-1111-1111-111111111111",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "br_book"`).
		WithArgs(book.ID, book.Title, book.Author, book.ISBN, book.Description, book.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(book)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
