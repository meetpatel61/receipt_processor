package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Receipt represents the structure of a receipt with various fields like retailer name, purchase date and time, total amount, and list of items.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

// Item represents an individual item on the receipt with a description and price.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// ReceiptStore is a thread-safe structure that holds receipts and their associated points in memory.
type ReceiptStore struct {
	sync.Mutex
	receipts map[string]int
}

// Initialize a new in-memory store to hold receipt data.
var store = &ReceiptStore{
	receipts: make(map[string]int),
}

func main() {
	// Define HTTP routes for the API endpoints
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	// Start the server on port 8080
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

// processReceiptHandler handles the /receipts/process endpoint, which accepts a JSON receipt,
// calculates points, generates an ID for the receipt, and stores it in memory.
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body into a Receipt struct
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt and calculate points based on the receipt details
	id := uuid.New().String()
	points := calculatePoints(receipt)

	// Store the receipt ID and its associated points in memory
	store.Lock()
	store.receipts[id] = points
	store.Unlock()

	// Return the generated ID as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// getPointsHandler handles the /receipts/{id}/points endpoint, which retrieves the points for a specific receipt ID.
func getPointsHandler(w http.ResponseWriter, r *http.Request) {

	// Extract the ID from the URL and remove the trailing /points suffix
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	// Retrieve the points for the given receipt ID from the in-memory store
	store.Lock()
	points, exists := store.receipts[id]
	store.Unlock()

	// If the ID is not found, return a 404 error
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Return the points as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

// calculatePoints calculates the total points for a given receipt based on predefined rules.
func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)

	// Parse the total as a float to apply rules based on the total amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if math.Mod(total, 1.0) == 0 {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Additional points based on the length of item descriptions
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)

		// If the length of the description is a multiple of 3, award points based on 20% of the item price
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the day of the purchase date is odd
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00 PM and 4:00 PM
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	// Return the total calculated points
	return points
}

// countAlphanumeric counts the number of alphanumeric characters in a given string.
// This is used to calculate points based on the retailer name.
func countAlphanumeric(s string) int {
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(re.FindAllString(s, -1))
}
