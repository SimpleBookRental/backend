package repositories

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func setupTokenMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestTokenRepository_FindTokenByValue_Found(t *testing.T) {
	gdb, mock, cleanup := setupTokenMockDB(t)
	defer cleanup()
	repo := NewTokenRepository(gdb)

	rows := sqlmock.NewRows([]string{"user_id", "token", "token_type", "expires_at"}).
		AddRow("11111111-1111-1111-1111-111111111111", "token-value", "access", time.Now().Add(time.Hour))

	mock.ExpectQuery(`SELECT .* FROM "br_issued_token"`).
		WithArgs("token-value", sqlmock.AnyArg()).
		WillReturnRows(rows)

	token, err := repo.FindTokenByValue("token-value")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "token-value", token.Token)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenRepository_FindTokenByValue_NotFound(t *testing.T) {
	gdb, mock, cleanup := setupTokenMockDB(t)
	defer cleanup()
	repo := NewTokenRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_issued_token"`).
		WithArgs("token-value", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "token", "token_type", "expires_at"}))

	token, err := repo.FindTokenByValue("token-value")
	assert.NoError(t, err)
	assert.Nil(t, token)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenRepository_FindTokenByValue_Error(t *testing.T) {
	gdb, mock, cleanup := setupTokenMockDB(t)
	defer cleanup()
	repo := NewTokenRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_issued_token"`).
		WithArgs("token-value", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)

	token, err := repo.FindTokenByValue("token-value")
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenRepository_CreateToken(t *testing.T) {
	gdb, mock, cleanup := setupTokenMockDB(t)
	defer cleanup()
	repo := NewTokenRepository(gdb)

	token := &models.IssuedToken{
		UserID:    "11111111-1111-1111-1111-111111111111",
		Token:     "token-value",
		TokenType: "access",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "br_issued_token"`).
		WithArgs(sqlmock.AnyArg(), token.UserID, token.Token, token.TokenType, token.ExpiresAt, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateToken(token)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
