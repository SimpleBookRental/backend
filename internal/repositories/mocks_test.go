package repositories

import (
	"log"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SimpleBookRental/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMockDB creates a new mock database connection
func NewMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Printf("An error '%s' was not expected when opening a stub database connection", err)
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Printf("An error '%s' was not expected when opening a gorm database connection", err)
		return nil, nil, err
	}

	return gormDB, mock, nil
}

// MockUser creates a mock user for testing
func MockUser() *models.User {
	return &models.User{
		ID:        "test-user-id",
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// MockBook creates a mock book for testing
func MockBook() *models.Book {
	return &models.Book{
		ID:          "test-book-id",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "test-user-id",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// MockIssuedToken creates a mock issued token for testing
func MockIssuedToken() *models.IssuedToken {
	return &models.IssuedToken{
		ID:        "test-token-id",
		UserID:    "test-user-id",
		Token:     "test-token-value",
		TokenType: string(models.AccessToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AnyTime is a matcher for any time value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

// ExpectUserQueries sets up expectations for user queries
func ExpectUserQueries(mock sqlmock.Sqlmock, user *models.User) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "br_user" WHERE id = $1 ORDER BY "br_user"."id" LIMIT 1`)).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
			AddRow(user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt))
}

// ExpectBookQueries sets up expectations for book queries
func ExpectBookQueries(mock sqlmock.Sqlmock, book *models.Book) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "br_book" WHERE id = $1 ORDER BY "br_book"."id" LIMIT 1`)).
		WithArgs(book.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "isbn", "description", "user_id", "created_at", "updated_at"}).
			AddRow(book.ID, book.Title, book.Author, book.ISBN, book.Description, book.UserID, book.CreatedAt, book.UpdatedAt))
}

// ExpectTokenQueries sets up expectations for token queries
func ExpectTokenQueries(mock sqlmock.Sqlmock, token *models.IssuedToken) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "br_issued_token" WHERE token = $1 ORDER BY "br_issued_token"."id" LIMIT 1`)).
		WithArgs(token.Token).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "token_type", "expires_at", "is_revoked", "revoked_at", "created_at", "updated_at"}).
			AddRow(token.ID, token.UserID, token.Token, token.TokenType, token.ExpiresAt, token.IsRevoked, token.RevokedAt, token.CreatedAt, token.UpdatedAt))
}
