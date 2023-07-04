# goapi

This is a Go API application that provides [brief description of the purpose or functionality of your API].

Installation
============= 

1)Clone the repository:

   - bash

    - git clone https://github.com/Carol757/goapi.git

2)Navigate to the project directory:

    -cd goapi

3)Install the required dependencies:

    - go

    - go mod download

Usage
============= 

1)Run the application:

    - go

    - go run main.go

    The API server will start on http://localhost:8080.


2)Use a tool like cURL or Postman to interact with the API endpoints.

API Endpoints
---------------

- Path: /receipts/{id}/points

- Method: GET

- Response: A JSON object containing the number of points awarded.

- Description: A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.]

Example:

bash

curl -X GET http://localhost:8080/receipts/{id}/points

Response:

{ "points": 32 }


Path: /receipts/process

Method: POST

Payload: Receipt JSON

Response: JSON containing an id for the receipt.

-Description: Takes in a JSON receipt (see example in the example directory) and returns a JSON object with an ID generated by your code.

Example:

bash

curl -X POST -H "Content-Type: application/json" -d '{"key": "value"}' http://localhost:8080/receipts/process

Response:

{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
