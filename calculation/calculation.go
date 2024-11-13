package calculation

import (
	"fetch-app/server" // Corrected import path for Receipt
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Helper function to calculate points
func CalculatePoints(receipt server.Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)

	// Rule 2: 50 points if the total is a round dollar amount (no cents)
	if isRoundDollar(receipt.Total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if isMultipleOfQuarter(receipt.Total) {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Points based on item descriptions
	for _, item := range receipt.Items {
		points += pointsForItemDescription(item)
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	if isOddDay(receipt.PurchaseDate.String()) {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if isBetweenTwoAndFourPM(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

// Rule 1: Count alphanumeric characters in retailer name
func countAlphanumeric(s string) int {
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(re.FindAllString(s, -1))
}

// Rule 2: Check if total is a round dollar amount (i.e., no cents)
func isRoundDollar(total string) bool {
	// Try to parse the total as a float
	val, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return false
	}
	return val == math.Floor(val)
}

// Rule 3: Check if total is a multiple of 0.25
func isMultipleOfQuarter(total string) bool {
	val, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return false
	}
	// Check if it's a multiple of 0.25
	return math.Mod(val, 0.25) == 0
}

// Rule 4: 5 points for every two items on the receipt
// This is handled directly by dividing the length of the items array

// Rule 5: Points based on item descriptions
func pointsForItemDescription(item server.Item) int {
	// Trim the description (remove leading and trailing spaces)
	trimmedDesc := strings.TrimSpace(item.ShortDescription)
	// Check if length is a multiple of 3
	if len(trimmedDesc)%3 == 0 {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return 0
		}
		// Multiply price by 0.2 and round up
		return int(math.Ceil(price * 0.2))
	}
	return 0
}

// Rule 6: 6 points if the day in the purchase date is odd
func isOddDay(date string) bool {
	// Parse the date (assume format "YYYY-MM-DD")
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	// Check if the day of the month is odd
	return parsedDate.Day()%2 != 0
}

// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
func isBetweenTwoAndFourPM(timeStr string) bool {
	// Parse the time (assume 24-hour format "HH:MM")
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	// Check if the time is between 2:00 PM and 4:00 PM
	return parsedTime.Hour() >= 14 && parsedTime.Hour() < 16
}
