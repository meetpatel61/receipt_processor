Receipt Processor

This project is a simple Receipt Processor web service written in Go. The service provides a REST API for submitting receipts, calculating points based on specific rules, and retrieving the awarded points.

Project Structure

	•	main.go: Contains the main application code, including API handlers and point calculation logic.
	•	go.mod: Defines the module, Go version, and dependencies.
	•	go.sum: Verifies the integrity of dependencies with cryptographic hashes.

Prerequisites

	•	Go version 1.23.2 or higher installed on your system.

Getting Started

1. Clone the Repository

git clone https://github.com/your-username/receipt_processor.git
cd receipt_processor

2. Install Dependencies

Initialize and install the required dependencies:

go mod tidy

3. Run the Application

To start the server, run:

go run main.go

The server will start on port 8080. You should see a message like:

Starting server on :8080

API Documentation

The API provides two endpoints:

1. Process Receipt

	•	Endpoint: /receipts/process
	•	Method: POST
	•	Description: Submits a receipt for processing. The server calculates the points based on the receipt’s content and returns a unique receipt ID.
	•	Request Payload: JSON object representing the receipt.

{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {
            "shortDescription": "Pepsi - 12-oz",
            "price": "1.25"
        }
    ]
}


	•	Response: JSON object with a unique id for the receipt.

{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }


	•	Example Request:

curl -X POST -H "Content-Type: application/json" -d @simple-receipt.json http://localhost:8080/receipts/process



2. Get Points

	•	Endpoint: /receipts/{id}/points
	•	Method: GET
	•	Description: Retrieves the points awarded for a specific receipt.
	•	Path Parameter: id - The unique ID of the receipt (obtained from the /receipts/process endpoint).
	•	Response: JSON object containing the points awarded.

{ "points": 32 }


	•	Example Request:

curl http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points



Points Calculation Rules

Points are awarded based on the following rules:
	1.	Retailer Name: One point for every alphanumeric character in the retailer name.
	2.	Round Dollar Total: 50 points if the total is a round dollar amount with no cents.
	3.	Total Multiple of 0.25: 25 points if the total is a multiple of 0.25.
	4.	Item Count: 5 points for every two items on the receipt.
	5.	Item Description Length: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points awarded for that item.
	6.	Odd Purchase Day: 6 points if the purchase day is an odd number.
	7.	Purchase Time: 10 points if the time of purchase is between 2:00 PM and 4:00 PM.

Example Usage

	1.	Start the server:

go run main.go


	2.	Process a receipt (using simple-receipt.json as an example):

curl -X POST -H "Content-Type: application/json" -d @simple-receipt.json http://localhost:8080/receipts/process

	•	Response:

{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }


	3.	Get points for the receipt ID:

curl http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points

	•	Response:

{ "points": 32 }



Development Notes

	•	Data Persistence: This application uses in-memory storage (a Go map) to store receipt data. All data is lost when the application stops, as no database is used.
	•	Dependencies:
	•	github.com/google/uuid - Used to generate unique IDs for each receipt.

License

This project is open-source and available under the MIT License.
