package main

import (
	"bytes"
	"encoding/json"
	"fetch-app/calculation"
	"fetch-app/server"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestPostReceiptsProcess tests the PostReceiptsProcess handler.
func TestPostReceiptsProcess(t *testing.T) {
	e := echo.New()
	handler := &ReceiptHandler{}

	// Create a test request with a valid receipt
	receipt := server.PostReceiptsProcessJSONRequestBody{
		Retailer:     "M&M Corner Market",
		PurchaseDate: types.Date{Time: time.Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC)},
		PurchaseTime: "14:33",
		Items: []server.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}

	// Convert the receipt to JSON
	reqBody, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBufferString(string(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Call the handler
	e.POST("/receipts/process", handler.PostReceiptsProcess)
	e.ServeHTTP(rec, req)

	// Check the response status
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the ID is returned in the response
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "id")

	// Verify the receipt was added to the storage
	receiptID := response["id"]
	_, exists := receiptStorage.Receipts[receiptID]
	assert.True(t, exists)
}

// TestGetReceiptsIdPoints tests the GetReceiptsIdPoints handler.
func TestGetReceiptsIdPoints(t *testing.T) {
	e := echo.New()
	handler := &ReceiptHandler{}

	// First, create a receipt and store it manually for testing
	receipt := server.PostReceiptsProcessJSONRequestBody{
		Retailer:     "M&M Corner Market",
		PurchaseDate: types.Date{Time: time.Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC)},
		PurchaseTime: "14:33",
		Items: []server.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}
	receiptID := uuid.New().String() // Generate a new receipt ID
	receiptStorage.Receipts[receiptID] = receipt

	// Create a test request to retrieve points for the stored receipt
	req := httptest.NewRequest(http.MethodGet, "/receipts/"+receiptID+"/points", nil)
	rec := httptest.NewRecorder()

	// Call the handler
	e.GET("/receipts/:id/points", func(c echo.Context) error {
		return handler.GetReceiptsIdPoints(c, receiptID) // Call the handler with the ID
	})
	e.ServeHTTP(rec, req)

	// Check the response status
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the points are returned in the response
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "points")
	assert.IsType(t, float64(0), response["points"])

	// Check the points logic (modify based on your actual calculation)
	// If expectedPoints is your calculated points, you can assert like this:
	expectedPoints := calculation.CalculatePoints(receipt)
	assert.Equal(t, float64(expectedPoints), response["points"])
}

// TestGetReceiptsIdPointsNotFound tests the case where the receipt does not exist.
func TestGetReceiptsIdPointsNotFound(t *testing.T) {
	e := echo.New()
	handler := &ReceiptHandler{}

	// Create a test request to retrieve points for a non-existing receipt
	nonExistentID := uuid.New().String() // Random ID for testing
	req := httptest.NewRequest(http.MethodGet, "/receipts/"+nonExistentID+"/points", nil)
	rec := httptest.NewRecorder()

	// Call the handler
	e.GET("/receipts/:id/points", func(c echo.Context) error {
		return handler.GetReceiptsIdPoints(c, nonExistentID) // Call the handler with the ID
	})
	e.ServeHTTP(rec, req)

	// Check the response status
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Verify the message in the response
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Equal(t, fmt.Sprintf("Receipt with ID %s not found", nonExistentID), response["message"])
}
