// Unit tests for BookController using testify and generated mocks.
package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupGinBook() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestBookController_Create_Success_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.POST("/books", func(ctx *gin.Context) {
		ctx.Set("user_id", "11111111-1111-1111-1111-111111111111")
		ctx.Set("role", models.UserRole)
		controller.Create(ctx)
	})

	bookCreate := models.BookCreate{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
	}
	expectedBook := &models.Book{
		ID:     "book-1",
		Title:  "Book 1",
		Author: "Author 1",
		UserID: "11111111-1111-1111-1111-111111111111",
	}
	mockService.EXPECT().Create(gomock.Any()).Return(expectedBook, nil)

	body, _ := json.Marshal(bookCreate)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestBookController_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.POST("/books", func(ctx *gin.Context) {
		ctx.Set("user_id", "11111111-1111-1111-1111-111111111111")
		ctx.Set("role", models.UserRole)
		controller.Create(ctx)
	})

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookController_Create_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.POST("/books", func(ctx *gin.Context) {
		ctx.Set("user_id", "11111111-1111-1111-1111-111111111111")
		ctx.Set("role", models.UserRole)
		controller.Create(ctx)
	})

	bookCreate := models.BookCreate{
		Title:  "Book 1",
		Author: "Author 1",
		ISBN:   "1234567890",
	}
	mockService.EXPECT().Create(gomock.Any()).Return(nil, errors.New("service error"))

	body, _ := json.Marshal(bookCreate)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookController_GetByID_Success_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.GET("/books/:id", func(ctx *gin.Context) {
		ctx.Set("user_id", "user-1")
		ctx.Set("role", models.UserRole)
		controller.GetByID(ctx)
	})

	expectedBook := &models.Book{
		ID:     "book-1",
		Title:  "Book 1",
		Author: "Author 1",
		UserID: "user-1",
	}
	mockService.EXPECT().GetByID("book-1").Return(expectedBook, nil)

	req, _ := http.NewRequest("GET", "/books/book-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookController_GetByID_Forbidden_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.GET("/books/:id", func(ctx *gin.Context) {
		ctx.Set("user_id", "user-1")
		ctx.Set("role", models.UserRole)
		controller.GetByID(ctx)
	})

	// Book belongs to another user
	book := &models.Book{
		ID:     "book-1",
		Title:  "Book 1",
		Author: "Author 1",
		UserID: "user-2",
	}
	mockService.EXPECT().GetByID("book-1").Return(book, nil)

	req, _ := http.NewRequest("GET", "/books/book-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBookController_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.GET("/books/:id", func(ctx *gin.Context) {
		controller.GetByID(ctx)
	})

	mockService.EXPECT().GetByID("book-404").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/books/book-404", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookController_GetAll_AdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.GET("/books", func(ctx *gin.Context) {
		ctx.Set("user_id", "admin-1")
		ctx.Set("role", models.AdminRole)
		controller.GetAll(ctx)
	})

	books := []models.Book{
		{ID: "book-1", UserID: "user-1"},
		{ID: "book-2", UserID: "user-2"},
	}
	mockService.EXPECT().GetAll().Return(books, nil)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookController_GetAll_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.GET("/books", func(ctx *gin.Context) {
		ctx.Set("user_id", "user-1")
		ctx.Set("role", models.UserRole)
		controller.GetAll(ctx)
	})

	books := []models.Book{
		{ID: "book-1", UserID: "user-1"},
	}
	mockService.EXPECT().GetByUserID("user-1").Return(books, nil)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookController_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.GET("/books", func(ctx *gin.Context) {
		controller.GetAll(ctx)
	})

	mockService.EXPECT().GetAll().Return(nil, errors.New("db error"))

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBookController_Update_Success_AdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.PUT("/books/:id", func(ctx *gin.Context) {
		ctx.Set("user_id", "admin-1")
		ctx.Set("role", models.AdminRole)
		controller.Update(ctx)
	})

	existingBook := &models.Book{
		ID:     "book-1",
		UserID: "user-1",
	}
	updatedBook := &models.Book{
		ID:     "book-1",
		UserID: "user-1",
		Title:  "Updated",
	}
	mockService.EXPECT().GetByID("book-1").Return(existingBook, nil)
	mockService.EXPECT().Update("book-1", gomock.Any()).Return(updatedBook, nil)

	update := models.BookUpdate{Title: "Updated"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/books/book-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookController_Update_Forbidden_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.PUT("/books/:id", func(ctx *gin.Context) {
		ctx.Set("user_id", "user-1")
		ctx.Set("role", models.UserRole)
		controller.Update(ctx)
	})

	existingBook := &models.Book{
		ID:     "book-1",
		UserID: "user-2",
	}
	mockService.EXPECT().GetByID("book-1").Return(existingBook, nil)

	update := models.BookUpdate{Title: "Updated"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/books/book-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBookController_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.PUT("/books/:id", func(ctx *gin.Context) {
		controller.Update(ctx)
	})

	mockService.EXPECT().GetByID("book-404").Return(nil, errors.New("not found"))

	update := models.BookUpdate{Title: "Updated"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/books/book-404", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBookController_Delete_Success_AdminRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.DELETE("/books/:id", func(ctx *gin.Context) {
		ctx.Set("user_id", "admin-1")
		ctx.Set("role", models.AdminRole)
		controller.Delete(ctx)
	})

	existingBook := &models.Book{
		ID:     "book-1",
		UserID: "user-1",
	}
	mockService.EXPECT().GetByID("book-1").Return(existingBook, nil)
	mockService.EXPECT().Delete("book-1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/books/book-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookController_Delete_Forbidden_UserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.DELETE("/books/:id", func(ctx *gin.Context) {
		ctx.Set("user_id", "user-1")
		ctx.Set("role", models.UserRole)
		controller.Delete(ctx)
	})

	existingBook := &models.Book{
		ID:     "book-1",
		UserID: "user-2",
	}
	mockService.EXPECT().GetByID("book-1").Return(existingBook, nil)

	req, _ := http.NewRequest("DELETE", "/books/book-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBookController_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := NewBookController(mockService)
	router := setupGinBook()
	router.DELETE("/books/:id", func(ctx *gin.Context) {
		controller.Delete(ctx)
	})

	mockService.EXPECT().GetByID("book-404").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("DELETE", "/books/book-404", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
