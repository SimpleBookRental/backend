// Unit tests for BookService using gomock and testify.
package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
)

func validUUID() string {
	return "11111111-1111-1111-1111-111111111111"
}

func TestBookService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
		UserID: validUUID(),
	}
	user := &models.User{ID: validUUID()}
	mockUserRepo.EXPECT().FindByID(bookCreate.UserID).Return(user, nil)
	mockBookRepo.EXPECT().FindByISBN(bookCreate.ISBN).Return(nil, nil)
	mockBookRepo.EXPECT().Create(gomock.Any()).Return(nil)

	book, err := service.Create(bookCreate)
	assert.NoError(t, err)
	assert.Equal(t, bookCreate.Title, book.Title)
	assert.Equal(t, bookCreate.UserID, book.UserID)
}

func TestBookService_Create_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
		UserID: validUUID(),
	}
	mockUserRepo.EXPECT().FindByID(bookCreate.UserID).Return(nil, nil)

	book, err := service.Create(bookCreate)
	assert.Nil(t, book)
	assert.ErrorContains(t, err, "user not found")
}

func TestBookService_Create_ISBNExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
		UserID: validUUID(),
	}
	user := &models.User{ID: validUUID()}
	mockUserRepo.EXPECT().FindByID(bookCreate.UserID).Return(user, nil)
	mockBookRepo.EXPECT().FindByISBN(bookCreate.ISBN).Return(&models.Book{}, nil)

	book, err := service.Create(bookCreate)
	assert.Nil(t, book)
	assert.ErrorContains(t, err, "ISBN already exists")
}

func TestBookService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
		UserID: validUUID(),
	}
	mockUserRepo.EXPECT().FindByID(bookCreate.UserID).Return(nil, errors.New("db error"))

	book, err := service.Create(bookCreate)
	assert.Nil(t, book)
	assert.ErrorContains(t, err, "error checking user")
}

func TestBookService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	book := &models.Book{ID: validUUID()}
	mockBookRepo.EXPECT().FindByID(book.ID).Return(book, nil)

	got, err := service.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book.ID, got.ID)
}

func TestBookService_GetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	got, err := service.GetByID("invalid")
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "invalid book ID")
}

func TestBookService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	mockBookRepo.EXPECT().FindByID(validUUID()).Return(nil, nil)

	got, err := service.GetByID(validUUID())
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "book not found")
}

func TestBookService_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	books := []models.Book{{ID: "1"}, {ID: "2"}}
	mockBookRepo.EXPECT().FindAll().Return(books, nil)

	got, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, books, got)
}

func TestBookService_GetByUserID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	user := &models.User{ID: validUUID()}
	books := []models.Book{{ID: "1"}, {ID: "2"}}
	mockUserRepo.EXPECT().FindByID(user.ID).Return(user, nil)
	mockBookRepo.EXPECT().FindByUserID(user.ID).Return(books, nil)

	got, err := service.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, books, got)
}

func TestBookService_GetByUserID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	got, err := service.GetByUserID("invalid")
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "invalid user ID")
}

func TestBookService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	book := &models.Book{ID: validUUID(), Title: "Old"}
	bookUpdate := &models.BookUpdate{Title: "New"}
	mockBookRepo.EXPECT().FindByID(book.ID).Return(book, nil)
	mockBookRepo.EXPECT().Update(book).Return(nil)

	got, err := service.Update(book.ID, bookUpdate)
	assert.NoError(t, err)
	assert.Equal(t, "New", got.Title)
}

func TestBookService_Update_ISBNExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	book := &models.Book{ID: validUUID(), ISBN: "old"}
	bookUpdate := &models.BookUpdate{ISBN: "new"}
	mockBookRepo.EXPECT().FindByID(book.ID).Return(book, nil)
	mockBookRepo.EXPECT().FindByISBN("new").Return(&models.Book{}, nil)

	got, err := service.Update(book.ID, bookUpdate)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "ISBN already exists")
}

func TestBookService_Update_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	book := &models.Book{ID: validUUID()}
	bookUpdate := &models.BookUpdate{UserID: validUUID()}
	mockBookRepo.EXPECT().FindByID(book.ID).Return(book, nil)
	mockUserRepo.EXPECT().FindByID(bookUpdate.UserID).Return(nil, nil)

	got, err := service.Update(book.ID, bookUpdate)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "user not found")
}

func TestBookService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	book := &models.Book{ID: validUUID()}
	mockBookRepo.EXPECT().FindByID(book.ID).Return(book, nil)
	mockBookRepo.EXPECT().Delete(book.ID).Return(nil)

	err := service.Delete(book.ID)
	assert.NoError(t, err)
}

func TestBookService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	mockBookRepo.EXPECT().FindByID(validUUID()).Return(nil, nil)

	err := service.Delete(validUUID())
	assert.ErrorContains(t, err, "book not found")
}

func TestBookService_Delete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewBookService(mockBookRepo, mockUserRepo)

	err := service.Delete("invalid")
	assert.ErrorContains(t, err, "invalid book ID")
}
