# Receipt Points Calculator

This project provides an API to calculate points based on receipt data. You can add receipts and query the calculated points using simple HTTP requests.

# Testing the application
All tests can be run using the following:
   - `go test fetch-app`

# Running the Application
### Build and Run the Application Using Docker
To build and run the application in Docker, follow these steps:

1. Clone the repository and navigate to the project folder.

2. Build the Docker image:
   - `docker build -t fetch-app .`

3. Run the application:
   - `docker run fetch-app`

This will start the application on `localhost:8080`.

# Interacting with the API
Once the application is running, you can interact with it using curl commands from the command line.

### Add a Receipt
Use a POST request to add a receipt for processing. Replace the example data in the curl command with the actual receipt data.

Example curl command to submit a receipt:


```bash
curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d "{\"retailer\":\"M^&M Corner Market\",\"purchaseDate\":\"2022-03-20\",\"purchaseTime\":\"14:33\",\"items\":[{\"shortDescription\":\"Gatorade\",\"price\":\"2.25\"},{\"shortDescription\":\"Gatorade\",\"price\":\"2.25\"},{\"shortDescription\":\"Gatorade\",\"price\":\"2.25\"},{\"shortDescription\":\"Gatorade\",\"price\":\"2.25\"}],\"total\":\"9.00\"}"
```
This will return a unique ID associated with the receipt. For example:

```bash
{
  "id": "2b2d8024-acb6-4eaa-9ed4-dcae58dd0331"
}
```

### Get Points for a Receipt
Once you have the receipt ID, you can query the points for the receipt using the following GET request:

```bash
curl -X GET http://localhost:8080/receipts/2b2d8024-acb6-4eaa-9ed4-dcae58dd0331/points
```

The server will respond with the calculated points, for example:
```
{
  "points": 109
}
```

### Example Receipt Data
Here is an example of a receipt that you can use with the above curl commands:

```bash
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```
