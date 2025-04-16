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

func setupBookController(t *testing.T) (*gin.Engine, *mocks.MockBookServiceInterface) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockBookServiceInterface(ctrl)
	controller := &BookController{
		bookService: mockService,
	}

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/books", controller.Create)
		v1.GET("/books", controller.GetAll)
		v1.GET("/books/:id", controller.GetByID)
		v1.PUT("/books/:id", controller.Update)
		v1.DELETE("/books/:id", controller.Delete)
	}

	return router, mockService
}

func TestBookController_Create(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookCreate := &models.BookCreate{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	book := &models.Book{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		Title:       bookCreate.Title,
		Author:      bookCreate.Author,
		ISBN:        bookCreate.ISBN,
		Description: bookCreate.Description,
		UserID:      bookCreate.UserID,
	}

	// Expectations
	mockService.EXPECT().Create(gomock.AssignableToTypeOf(&models.BookCreate{})).Return(book, nil)

	// Create request
	body, _ := json.Marshal(bookCreate)
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Book created successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Create_InvalidRequest(t *testing.T) {
	// Setup
	router, _ := setupBookController(t)

	// Invalid request (missing required fields)
	bookCreate := map[string]interface{}{
		"title": "Test Book",
		// Missing author, isbn, and user_id
	}

	// Create request
	body, _ := json.Marshal(bookCreate)
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Invalid request body", response["message"])

	// No expectations to verify as the service should not be called
}

func TestBookController_Create_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookCreate := &models.BookCreate{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	// Expectations
	mockService.EXPECT().Create(gomock.AssignableToTypeOf(&models.BookCreate{})).Return(nil, errors.New("service error"))

	// Create request
	body, _ := json.Marshal(bookCreate)
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to create book", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_GetAll(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	books := []models.Book{
		{
			ID:     "book-1",
			Title:  "Book 1",
			Author: "Author 1",
		},
		{
			ID:     "book-2",
			Title:  "Book 2",
			Author: "Author 2",
		},
	}

	// Expectations
	mockService.EXPECT().GetAll().Return(books, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/books", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Books retrieved successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_GetAll_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	// Expectations
	mockService.EXPECT().GetAll().Return([]models.Book{}, errors.New("service error"))

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/books", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to retrieve books", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_GetByID(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	book := &models.Book{
		ID:     bookID,
		Title:  "Test Book",
		Author: "Test Author",
	}

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(book, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/books/"+bookID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Book retrieved successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_GetByID_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(nil, errors.New("service error"))

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/books/"+bookID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "service error", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Update(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	bookUpdate := &models.BookUpdate{
		Title:       "Updated Book",
		Author:      "Updated Author",
		ISBN:        "ISBN-456",
		Description: "Updated description",
	}

	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Original Book",
		Author:      "Original Author",
		ISBN:        "ISBN-123",
		Description: "Original description",
		UserID:      "user-id",
	}

	updatedBook := &models.Book{
		ID:          bookID,
		Title:       bookUpdate.Title,
		Author:      bookUpdate.Author,
		ISBN:        bookUpdate.ISBN,
		Description: bookUpdate.Description,
		UserID:      "user-id",
	}

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(existingBook, nil)
	mockService.EXPECT().Update(bookID, gomock.AssignableToTypeOf(&models.BookUpdate{})).Return(updatedBook, nil)

	// Create request
	body, _ := json.Marshal(bookUpdate)
	req, _ := http.NewRequest("PUT", "/api/v1/books/"+bookID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Book updated successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Update_InvalidRequest(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Original Book",
		Author:      "Original Author",
		ISBN:        "ISBN-123",
		Description: "Original description",
		UserID:      "user-id",
	}

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(existingBook, nil)

	// Invalid request (invalid JSON)
	req, _ := http.NewRequest("PUT", "/api/v1/books/"+bookID, bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Invalid request body", response["message"])

	// No expectations to verify as the service should not be called
}

func TestBookController_Update_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	bookUpdate := &models.BookUpdate{
		Title:       "Updated Book",
		Author:      "Updated Author",
		ISBN:        "ISBN-456",
		Description: "Updated description",
	}

	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Original Book",
		Author:      "Original Author",
		ISBN:        "ISBN-123",
		Description: "Original description",
		UserID:      "user-id",
	}

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(existingBook, nil)
	mockService.EXPECT().Update(bookID, gomock.AssignableToTypeOf(&models.BookUpdate{})).Return(nil, errors.New("service error"))

	// Create request
	body, _ := json.Marshal(bookUpdate)
	req, _ := http.NewRequest("PUT", "/api/v1/books/"+bookID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to update book", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Delete(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Original Book",
		Author:      "Original Author",
		ISBN:        "ISBN-123",
		Description: "Original description",
		UserID:      "user-id",
	}

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(existingBook, nil)
	mockService.EXPECT().Delete(bookID).Return(nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/books/"+bookID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Book deleted successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Delete_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Original Book",
		Author:      "Original Author",
		ISBN:        "ISBN-123",
		Description: "Original description",
		UserID:      "user-id",
	}

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(existingBook, nil)
	mockService.EXPECT().Delete(bookID).Return(errors.New("service error"))

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/books/"+bookID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to delete book", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_GetByID_NotFound(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(nil, errors.New("book not found"))

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/books/"+bookID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "book not found", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Update_BookNotFound(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(nil, errors.New("book not found"))

	// Create request
	req, _ := http.NewRequest("PUT", "/api/v1/books/"+bookID, bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "book not found", response["message"])

	// Verify expectations handled by gomock controller
}

func TestBookController_Delete_BookNotFound(t *testing.T) {
	// Setup
	router, mockService := setupBookController(t)

	bookID := "123e4567-e89b-12d3-a456-426614174000"

	// Expectations
	mockService.EXPECT().GetByID(bookID).Return(nil, errors.New("book not found"))

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/books/"+bookID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "book not found", response["message"])

	// Verify expectations handled by gomock controller
}
