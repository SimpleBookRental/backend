// Unit tests for BookUserService using gomock and testify.
package services

import (
	"testing"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestBookUserService_TransferBookOwnership_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	bookID := "book-1"
	fromUserID := "user-1"
	toUserID := "user-2"
	book := &models.Book{ID: bookID, UserID: fromUserID}
	toUser := &models.User{ID: toUserID}

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			// Simulate transaction
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockBookRepoTx.EXPECT().FindByID(bookID).Return(book, nil)
	mockUserRepoTx.EXPECT().FindByID(toUserID).Return(toUser, nil)
	mockBookRepoTx.EXPECT().Update(book).Return(nil)

	err := service.TransferBookOwnership(bookID, fromUserID, toUserID)
	assert.NoError(t, err)
	assert.Equal(t, toUserID, book.UserID)
}

func TestBookUserService_TransferBookOwnership_BookNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	bookID := "book-1"
	fromUserID := "user-1"
	toUserID := "user-2"

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockBookRepoTx.EXPECT().FindByID(bookID).Return(nil, nil)

	err := service.TransferBookOwnership(bookID, fromUserID, toUserID)
	assert.ErrorContains(t, err, "book not found")
}

func TestBookUserService_TransferBookOwnership_NotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	bookID := "book-1"
	fromUserID := "user-1"
	toUserID := "user-2"
	book := &models.Book{ID: bookID, UserID: "other-user"}

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockBookRepoTx.EXPECT().FindByID(bookID).Return(book, nil)

	err := service.TransferBookOwnership(bookID, fromUserID, toUserID)
	assert.ErrorContains(t, err, "book does not belong to the specified user")
}

func TestBookUserService_TransferBookOwnership_TargetUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	bookID := "book-1"
	fromUserID := "user-1"
	toUserID := "user-2"
	book := &models.Book{ID: bookID, UserID: fromUserID}

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockBookRepoTx.EXPECT().FindByID(bookID).Return(book, nil)
	mockUserRepoTx.EXPECT().FindByID(toUserID).Return(nil, nil)

	err := service.TransferBookOwnership(bookID, fromUserID, toUserID)
	assert.ErrorContains(t, err, "target user not found")
}

func TestBookUserService_CreateBookWithUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	userID := "user-1"
	bookCreate := &models.BookCreateRequest{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
	}
	user := &models.User{ID: userID}

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockUserRepoTx.EXPECT().FindByID(userID).Return(user, nil)
	mockBookRepoTx.EXPECT().FindByISBN(bookCreate.ISBN).Return(nil, nil)
	mockBookRepoTx.EXPECT().Create(gomock.Any()).Return(nil)

	book, err := service.CreateBookWithUser(bookCreate, userID)
	assert.NoError(t, err)
	assert.Equal(t, bookCreate.Title, book.Title)
	assert.Equal(t, userID, book.UserID)
}

func TestBookUserService_CreateBookWithUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	userID := "user-1"
	bookCreate := &models.BookCreateRequest{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
	}

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockUserRepoTx.EXPECT().FindByID(userID).Return(nil, nil)

	book, err := service.CreateBookWithUser(bookCreate, userID)
	assert.Nil(t, book)
	assert.ErrorContains(t, err, "user not found")
}

func TestBookUserService_CreateBookWithUser_ISBNExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTxManager := mocks.NewMockTransactionManagerInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	service := NewBookUserService(mockTxManager, mockBookRepo, mockUserRepo)

	userID := "user-1"
	bookCreate := &models.BookCreateRequest{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
	}
	user := &models.User{ID: userID}

	mockTxManager.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(
		func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
	)
	mockBookRepoTx := mocks.NewMockBookRepositoryInterface(ctrl)
	mockUserRepoTx := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo.EXPECT().WithTx(gomock.Any()).Return(mockBookRepoTx, nil)
	mockUserRepo.EXPECT().WithTx(gomock.Any()).Return(mockUserRepoTx, nil)
	mockUserRepoTx.EXPECT().FindByID(userID).Return(user, nil)
	mockBookRepoTx.EXPECT().FindByISBN(bookCreate.ISBN).Return(&models.Book{}, nil)

	book, err := service.CreateBookWithUser(bookCreate, userID)
	assert.Nil(t, book)
	assert.ErrorContains(t, err, "ISBN already exists")
}
