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

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestUserRepository_Delete_Success(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "br_user"`).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete("11111111-1111-1111-1111-111111111111")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete_Error(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "br_user"`).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	err := repo.Delete("11111111-1111-1111-1111-111111111111")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_Success(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	user := &models.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Name:  "Updated User",
		Email: "updated@example.com",
		Password: "hashed",
		Role:  "USER",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "br_user"`).
		WithArgs(user.Name, user.Email, user.Password, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_Error(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	user := &models.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Name:  "Updated User",
		Email: "updated@example.com",
		Password: "hashed",
		Role:  "USER",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "br_user"`).
		WithArgs(user.Name, user.Email, user.Password, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID).
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	err := repo.Update(user)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindAll_Success(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow("11111111-1111-1111-1111-111111111111", "User1", "user1@example.com", "hashed", "USER").
		AddRow("22222222-2222-2222-2222-222222222222", "User2", "user2@example.com", "hashed", "USER")

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WillReturnRows(rows)

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "User1", users[0].Name)
	assert.Equal(t, "User2", users[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindAll_Empty(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}))

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindAll_Error(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WillReturnError(errors.New("db error"))

	users, err := repo.FindAll()
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_Found(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow("11111111-1111-1111-1111-111111111111", "Test User", "test@example.com", "hashed", "USER")

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("test@example.com", sqlmock.AnyArg()).
		WillReturnRows(rows)

	user, err := repo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("test@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}))

	user, err := repo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_Error(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("test@example.com", sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	user, err := repo.FindByEmail("test@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_Found(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow("11111111-1111-1111-1111-111111111111", "Test User", "test@example.com", "hashed", "USER")

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("11111111-1111-1111-1111-111111111111", sqlmock.AnyArg()).
		WillReturnRows(rows)
	// Mock preload Books
	mock.ExpectQuery(`SELECT .* FROM "br_book"`).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "user_id"}))

	user, err := repo.FindByID("11111111-1111-1111-1111-111111111111")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("11111111-1111-1111-1111-111111111111", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}))

	user, err := repo.FindByID("11111111-1111-1111-1111-111111111111")
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_Error(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	mock.ExpectQuery(`SELECT .* FROM "br_user"`).
		WithArgs("11111111-1111-1111-1111-111111111111", sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	user, err := repo.FindByID("11111111-1111-1111-1111-111111111111")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create(t *testing.T) {
	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewUserRepository(gdb)

	user := &models.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Name:  "Test User",
		Email: "test@example.com",
		Password: "hashed",
		Role:  "USER",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "br_user"`).
		WithArgs(user.ID, user.Name, user.Email, user.Password, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
