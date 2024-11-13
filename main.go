package main

import (
	"fetch-app/calculation"
	"fetch-app/server"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
)

// ReceiptStorage holds the receipts and provides a storage mechanism.
type ReceiptStorage struct {
	Receipts map[string]server.Receipt
}

// NewReceiptStorage initializes and returns a new ReceiptStorage instance.
func NewReceiptStorage() *ReceiptStorage {
	return &ReceiptStorage{
		Receipts: make(map[string]server.Receipt), // Initialize the map
	}
}

// Global receiptStorage instance that will be used across the application
var receiptStorage = NewReceiptStorage()

// ReceiptHandler implements the server routes related to receipt processing and points retrieval.
type ReceiptHandler struct{}

// PostReceiptsProcess handles the POST request to process a new receipt.
// It accepts a receipt in JSON format, stores it with a unique ID, and returns the ID in the response.
//
// Parameters:
//
//	ctx - The Echo context, which holds information about the request and response.
//
// Returns:
//
//	A JSON response containing the generated receipt ID if successful.
//	If the JSON is invalid or the binding fails, it returns a Bad Request (400) error with a relevant message.
func (h *ReceiptHandler) PostReceiptsProcess(ctx echo.Context) error {
	var receipt server.PostReceiptsProcessJSONRequestBody

	// Bind the incoming JSON request body to the receipt struct
	if err := ctx.Bind(&receipt); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
	}

	// Print the received receipt for debugging
	fmt.Printf("Received receipt: %+v\n", receipt)

	// Generate a unique ID for the receipt and store it in the receiptStorage map
	receiptID := uuid.New().String()
	receiptStorage.Receipts[receiptID] = receipt

	// Return a success response with the generated receipt ID
	return ctx.JSON(http.StatusOK, map[string]string{"id": receiptID})
}

// GetReceiptsIdPoints handles the GET request to retrieve points for a given receipt by ID.
// It checks if the receipt exists in storage, calculates points based on receipt data, and returns the result.
//
// Parameters:
//
//	ctx - The Echo context, which holds information about the request and response.
//	id  - The unique ID of the receipt whose points need to be retrieved.
//
// Returns:
//
//	A JSON response containing the calculated points if the receipt exists.
//	If the receipt does not exist, it returns a Not Found (404) error with a relevant message.
func (h *ReceiptHandler) GetReceiptsIdPoints(ctx echo.Context, id string) error {
	// Check if the receipt exists in the storage
	receipt, exists := receiptStorage.Receipts[id]
	if !exists {
		// If the receipt does not exist, return a 404 error with a relevant message
		return ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"message": fmt.Sprintf("Receipt with ID %s not found", id),
		})
	}

	// If the receipt exists, calculate and return the points
	points := calculation.CalculatePoints(receipt)

	// Return the points in the response
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"points": points,
	})
}

// main sets up the Echo server, registers the routes, and starts the application.
// It initializes the ReceiptHandler, sets up the routes, and begins listening on port 8080.
//
// Returns:
//
//	None (this is the entry point of the program, which starts the HTTP server).
func main() {
	// Create a new Echo instance
	e := echo.New()

	// Create the handler
	handler := &ReceiptHandler{}

	// Register the server routes
	server.RegisterHandlers(e, handler)

	// Start the Echo server on port 8080
	e.Start(":8080")
}
