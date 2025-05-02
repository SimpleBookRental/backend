package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestPaymentCRUD tests the payment API endpoints
func TestPaymentCRUD(t *testing.T) {
	// First create a test book
	createBookURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	bookData := map[string]interface{}{
		"title":       "Payment Test Book",
		"author":      "Payment Author",
		"publisher":   "Payment Publisher",
		"isbn":        "4444333322222",
		"categoryID":  1,
		"year":        2023,
		"description": "Book for payment test",
		"quantity":    3,
		"price":       14.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createBookURL, bookData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract book ID from response
	var createBookResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createBookResp); err != nil {
		t.Fatalf("Failed to decode create book response: %v", err)
	}
	
	bookData, ok := createBookResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from book response")
	}
	
	bookID, ok := bookData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract book ID from response")
	}
	
	// Create a rental that's overdue to generate a fee
	createRentalURL := fmt.Sprintf("%s/api/v1/rentals", baseURL)
	rentalData := map[string]interface{}{
		"bookID":  fmt.Sprintf("%.0f", bookID),
		"dueDate": time.Now().AddDate(0, 0, -5).Format("2006-01-02"), // 5 days ago
	}
	
	resp, err = makeAuthenticatedRequest("POST", createRentalURL, rentalData, memberToken)
	if err != nil {
		t.Fatalf("Failed to create rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract rental ID from response
	var createRentalResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createRentalResp); err != nil {
		t.Fatalf("Failed to decode create rental response: %v", err)
	}
	
	rentalData, ok = createRentalResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from rental response")
	}
	
	rentalID, ok := rentalData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract rental ID from response")
	}
	
	// Get fee amount
	feeURL := fmt.Sprintf("%s/api/v1/rentals/%s/fee", baseURL, fmt.Sprintf("%.0f", rentalID))
	resp, err = makeAuthenticatedRequest("GET", feeURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get fee: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var feeResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&feeResp); err != nil {
		t.Fatalf("Failed to decode fee response: %v", err)
	}
	
	feeData, ok := feeResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from fee response")
	}
	
	fee, ok := feeData["fee"].(float64)
	if !ok {
		t.Fatalf("Failed to extract fee from response")
	}
	
	// Create a payment
	createPaymentURL := fmt.Sprintf("%s/api/v1/payments", baseURL)
	paymentData := map[string]interface{}{
		"rentalID":     fmt.Sprintf("%.0f", rentalID),
		"amount":       fee,
		"paymentMethod": "credit_card",
		"reference":    "TEST-PAYMENT-REF-123",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createPaymentURL, paymentData, memberToken)
	if err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract payment ID from response
	var createPaymentResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createPaymentResp); err != nil {
		t.Fatalf("Failed to decode create payment response: %v", err)
	}
	
	paymentData, ok = createPaymentResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from payment response")
	}
	
	paymentID, ok := paymentData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract payment ID from response")
	}
	
	// Store payment ID for other tests
	testPaymentID = fmt.Sprintf("%.0f", paymentID)
	
	// Test get payment by ID
	getPaymentURL := fmt.Sprintf("%s/api/v1/payments/%s", baseURL, testPaymentID)
	resp, err = makeAuthenticatedRequest("GET", getPaymentURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get payment: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test list user payments
	listPaymentsURL := fmt.Sprintf("%s/api/v1/payments/user", baseURL)
	resp, err = makeAuthenticatedRequest("GET", listPaymentsURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to list user payments: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify returned rental status after payment
	getRentalURL := fmt.Sprintf("%s/api/v1/rentals/%s", baseURL, fmt.Sprintf("%.0f", rentalID))
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var getRentalResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&getRentalResp); err != nil {
		t.Fatalf("Failed to decode rental response: %v", err)
	}
	
	getRentalData, ok := getRentalResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from rental response")
	}
	
	// Since this varies by implementation, just check if we get any status
	rentalStatus, ok := getRentalData["status"].(string)
	if !ok || rentalStatus == "" {
		t.Errorf("Failed to get status from rental response")
	}
}

// TestPaymentValidation tests validation of payment data
func TestPaymentValidation(t *testing.T) {
	createPaymentURL := fmt.Sprintf("%s/api/v1/payments", baseURL)
	
	// Test invalid rental ID
	invalidRentalPayment := map[string]interface{}{
		"rentalID":     "999999", // Non-existent rental ID
		"amount":       10.0,
		"paymentMethod": "credit_card",
		"reference":    "TEST-INVALID-RENTAL",
	}
	
	resp, err := makeAuthenticatedRequest("POST", createPaymentURL, invalidRentalPayment, memberToken)
	if err != nil {
		t.Fatalf("Failed to make invalid rental payment request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusNotFound)
	
	// Test invalid amount (negative)
	invalidAmountPayment := map[string]interface{}{
		"rentalID":     testRentalID, // Using a valid rental ID from previous tests
		"amount":       -5.0,
		"paymentMethod": "credit_card",
		"reference":    "TEST-INVALID-AMOUNT",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createPaymentURL, invalidAmountPayment, memberToken)
	if err != nil {
		t.Fatalf("Failed to make invalid amount payment request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test invalid payment method
	invalidMethodPayment := map[string]interface{}{
		"rentalID":     testRentalID,
		"amount":       10.0,
		"paymentMethod": "invalid_method",
		"reference":    "TEST-INVALID-METHOD",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createPaymentURL, invalidMethodPayment, memberToken)
	if err != nil {
		t.Fatalf("Failed to make invalid method payment request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
}

// TestPaymentAdminOperations tests administrative payment operations
func TestPaymentAdminOperations(t *testing.T) {
	// Test admin can view all payments
	allPaymentsURL := fmt.Sprintf("%s/api/v1/payments/all", baseURL)
	
	// Admin should have access
	resp, err := makeAuthenticatedRequest("GET", allPaymentsURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access all payments as admin: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Librarian should have access
	resp, err = makeAuthenticatedRequest("GET", allPaymentsURL, nil, librianToken)
	if err != nil {
		t.Fatalf("Failed to access all payments as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Member should NOT have access
	resp, err = makeAuthenticatedRequest("GET", allPaymentsURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make request to admin endpoint as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test create payment by admin - admins can create payments on behalf of users
	if testRentalID != "" {
		adminCreatePaymentURL := fmt.Sprintf("%s/api/v1/payments/admin", baseURL)
		adminPaymentData := map[string]interface{}{
			"rentalID":      testRentalID,
			"amount":        15.0,
			"paymentMethod": "cash",
			"reference":     "ADMIN-CREATED-PAYMENT",
			"notes":         "Payment created by admin",
		}
		
		resp, err = makeAuthenticatedRequest("POST", adminCreatePaymentURL, adminPaymentData, adminToken)
		if err != nil {
			t.Fatalf("Failed to create admin payment: %v", err)
		}
		defer resp.Body.Close()
		
		checkStatusCode(t, resp, http.StatusCreated)
	}
}

// TestPaymentStats tests payment statistics endpoints
func TestPaymentStats(t *testing.T) {
	// Test payment statistics - admin only
	statsURL := fmt.Sprintf("%s/api/v1/payments/stats", baseURL)
	
	resp, err := makeAuthenticatedRequest("GET", statsURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get payment statistics: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var statsResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statsResp); err != nil {
		t.Fatalf("Failed to decode statistics response: %v", err)
	}
	
	statsData, ok := statsResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from statistics response")
	}
	
	// Check that we got statistics keys
	// Actual keys will depend on implementation, just check that we got something
	if len(statsData) == 0 {
		t.Errorf("Expected non-empty statistics data")
	}
	
	// Test payment statistics with date range
	rangeStatsURL := fmt.Sprintf("%s/api/v1/payments/stats?startDate=%s&endDate=%s", 
		baseURL, 
		time.Now().AddDate(0, -1, 0).Format("2006-01-02"),  // 1 month ago
		time.Now().Format("2006-01-02"))                    // today
	
	resp, err = makeAuthenticatedRequest("GET", rangeStatsURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get payment statistics with date range: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
}
