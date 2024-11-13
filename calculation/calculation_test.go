package calculation

import (
	"fetch-app/server"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Helper function to create a test receipt
func createTestReceipt() server.Receipt {
	return server.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: types.Date{Time: time.Date(2022, time.March, 20, 0, 0, 0, 0, time.UTC)}, // Using an "odd" date
		PurchaseTime: "14:33",                                                                 // Time between 2:00pm and 4:00pm
		Total:        "9.00",                                                                  // A valid round dollar amount
		Items: []server.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Candy", Price: "3.00"},
		},
	}
}

func TestCalculatePoints(t *testing.T) {
	receipt := createTestReceipt()

	// Expected points based on the rules
	expectedPoints := 0
	expectedPoints += countAlphanumeric(receipt.Retailer)        // Alphanumeric count for "M&M Corner Market"
	expectedPoints += 50                                         // Round dollar rule (9.00 is a round dollar)
	expectedPoints += 25                                         // Multiple of quarter rule (9.00 is a multiple of 0.25)
	expectedPoints += (len(receipt.Items) / 2) * 5               // Two items, so 5 points
	expectedPoints += pointsForItemDescription(receipt.Items[0]) // Gatorade has a description length of 8, which is a multiple of 3
	expectedPoints += pointsForItemDescription(receipt.Items[1]) // Another Gatorade item
	expectedPoints += pointsForItemDescription(receipt.Items[2]) // Candy has a length of 5, which is not a multiple of 3
	expectedPoints += 10                                         // Time between 2:00pm and 4:00pm

	actualPoints := CalculatePoints(receipt)

	assert.Equal(t, expectedPoints, actualPoints)
}

// Test for Rule 1: Alphanumeric count
func TestCountAlphanumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"M&M Corner Market", 14}, // "M&M Corner Market" has 15 alphanumeric characters
		{"123", 3},                // "123" has 3 alphanumeric characters
		{"!@#$%^", 0},             // No alphanumeric characters
		{"A B C", 3},              // "A B C" has 3 alphanumeric characters
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := countAlphanumeric(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

// Test for Rule 2: Round dollar check
func TestIsRoundDollar(t *testing.T) {
	tests := []struct {
		total    string
		expected bool
	}{
		{"10.00", true}, // Round dollar
		{"9.99", false}, // Not a round dollar
		{"5.50", false}, // Not a round dollar
		{"100", true},   // Round dollar
	}

	for _, test := range tests {
		t.Run(test.total, func(t *testing.T) {
			result := isRoundDollar(test.total)
			assert.Equal(t, test.expected, result)
		})
	}
}

// Test for Rule 3: Multiple of quarter check
func TestIsMultipleOfQuarter(t *testing.T) {
	tests := []struct {
		total    string
		expected bool
	}{
		{"9.00", true},  // Multiple of 0.25
		{"9.25", true},  // Multiple of 0.25
		{"9.10", false}, // Not a multiple of 0.25
		{"10.00", true}, // Multiple of 0.25
	}

	for _, test := range tests {
		t.Run(test.total, func(t *testing.T) {
			result := isMultipleOfQuarter(test.total)
			assert.Equal(t, test.expected, result)
		})
	}
}

// Test for Rule 5: Points based on item descriptions
func TestPointsForItemDescription(t *testing.T) {
	tests := []struct {
		item     server.Item
		expected int
	}{
		{server.Item{ShortDescription: "Gatorade", Price: "2.25"}, 0},  // Length of 8, not multiple of 3
		{server.Item{ShortDescription: "Coca-Cola", Price: "3.50"}, 1}, // Length of 9, multiple of 3
		{server.Item{ShortDescription: "Gum", Price: "1.50"}, 1},       // Length of 3, multiple of 3
		{server.Item{ShortDescription: "Candy", Price: "2.00"}, 0},     // Length of 5, not multiple of 3
	}

	for _, test := range tests {
		t.Run(test.item.ShortDescription, func(t *testing.T) {
			result := pointsForItemDescription(test.item)
			assert.Equal(t, test.expected, result)
		})
	}
}

// Test for Rule 6: Odd day check
func TestIsOddDay(t *testing.T) {
	tests := []struct {
		date     string
		expected bool
	}{
		{"2022-03-19", true},  // March 19th is odd
		{"2022-03-20", false}, // March 20th is even
		{"2022-03-21", true},  // March 21st is odd
	}

	for _, test := range tests {
		t.Run(test.date, func(t *testing.T) {
			result := isOddDay(test.date)
			assert.Equal(t, test.expected, result)
		})
	}
}

// Test for Rule 7: Time between 2:00pm and 4:00pm
func TestIsBetweenTwoAndFourPM(t *testing.T) {
	tests := []struct {
		timeStr  string
		expected bool
	}{
		{"14:30", true},  // Between 2:00 PM and 4:00 PM
		{"15:00", true},  // Exactly 3:00 PM
		{"16:00", false}, // Exactly 4:00 PM, not in the range
		{"13:00", false}, // Before 2:00 PM
	}

	for _, test := range tests {
		t.Run(test.timeStr, func(t *testing.T) {
			result := isBetweenTwoAndFourPM(test.timeStr)
			assert.Equal(t, test.expected, result)
		})
	}
}
