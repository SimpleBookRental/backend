// Unit tests for BookUserController using testify and generated mocks.
package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupGinBookUser() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestBookUserController_TransferBookOwnership_Success_AdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/:id/transfer", func(ctx *gin.Context) {
		ctx.Set("role", "ADMIN")
		controller.TransferBookOwnership(ctx)
	})

	reqBody := map[string]string{
		"from_user_id": "11111111-1111-1111-1111-111111111111",
		"to_user_id":   "22222222-2222-2222-2222-222222222222",
	}
	mockService.EXPECT().TransferBookOwnership("book-1", "11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222").Return(nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/books/book-1/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookUserController_TransferBookOwnership_Forbidden_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/:id/transfer", func(ctx *gin.Context) {
		ctx.Set("role", "USER")
		ctx.Set("user_id", "user-2")
		controller.TransferBookOwnership(ctx)
	})

	reqBody := map[string]string{
		"from_user_id": "11111111-1111-1111-1111-111111111111",
		"to_user_id":   "22222222-2222-2222-2222-222222222222",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/books/book-1/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBookUserController_TransferBookOwnership_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/:id/transfer", func(ctx *gin.Context) {
		ctx.Set("role", "ADMIN")
		controller.TransferBookOwnership(ctx)
	})

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/books/book-1/transfer", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookUserController_TransferBookOwnership_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/:id/transfer", func(ctx *gin.Context) {
		ctx.Set("role", "ADMIN")
		controller.TransferBookOwnership(ctx)
	})

	reqBody := map[string]string{
		"from_user_id": "11111111-1111-1111-1111-111111111111",
		"to_user_id":   "22222222-2222-2222-2222-222222222222",
	}
	mockService.EXPECT().TransferBookOwnership("book-1", "11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222").Return(errors.New("service error"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/books/book-1/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookUserController_CreateBookWithUser_Success_AdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/create_with_user", func(ctx *gin.Context) {
		ctx.Set("role", "ADMIN")
		controller.CreateBookWithUser(ctx)
	})

	reqBody := map[string]interface{}{
		"title":       "Book 1",
		"author":      "Author 1",
		"isbn":        "1234567890",
		"description": "desc",
		"user_id":     "11111111-1111-1111-1111-111111111111",
	}
	bookCreate := &models.BookCreateRequest{
		Title:       "Book 1",
		Author:      "Author 1",
		ISBN:        "1234567890",
		Description: "desc",
	}
	expectedBook := &models.Book{
		ID:     "book-1",
		Title:  "Book 1",
		Author: "Author 1",
		UserID: "11111111-1111-1111-1111-111111111111",
	}
	mockService.EXPECT().CreateBookWithUser(bookCreate, "11111111-1111-1111-1111-111111111111").Return(expectedBook, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/books/create_with_user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestBookUserController_CreateBookWithUser_Forbidden_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/create_with_user", func(ctx *gin.Context) {
		ctx.Set("role", "USER")
		ctx.Set("user_id", "user-2")
		controller.CreateBookWithUser(ctx)
	})

	reqBody := map[string]interface{}{
		"title":       "Book 1",
		"author":      "Author 1",
		"isbn":        "1234567890",
		"description": "desc",
		"user_id":     "11111111-1111-1111-1111-111111111111",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/books/create_with_user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBookUserController_CreateBookWithUser_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/create_with_user", func(ctx *gin.Context) {
		ctx.Set("role", "ADMIN")
		controller.CreateBookWithUser(ctx)
	})

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/books/create_with_user", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookUserController_CreateBookWithUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookUserServiceInterface(ctrl)
	controller := NewBookUserController(mockService)
	router := setupGinBookUser()
	router.POST("/books/create_with_user", func(ctx *gin.Context) {
		ctx.Set("role", "ADMIN")
		controller.CreateBookWithUser(ctx)
	})

	reqBody := map[string]interface{}{
		"title":       "Book 1",
		"author":      "Author 1",
		"isbn":        "1234567890",
		"description": "desc",
		"user_id":     "11111111-1111-1111-1111-111111111111",
	}
	bookCreate := &models.BookCreateRequest{
		Title:       "Book 1",
		Author:      "Author 1",
		ISBN:        "1234567890",
		Description: "desc",
	}
	mockService.EXPECT().CreateBookWithUser(bookCreate, "11111111-1111-1111-1111-111111111111").Return(nil, errors.New("service error"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/books/create_with_user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
