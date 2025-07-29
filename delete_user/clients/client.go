package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

// Define the request structure
type deleteUserRequest struct {
	ID int `json:"id"`
}

// Define the response structure
type deleteUserResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to NATS server
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// User ID to be deleted
	userID := 1 // Change this to the ID you want to delete

	// Create the request payload
	reqData := deleteUserRequest{ID: userID}
	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		log.Fatalf("Failed to marshal request data: %v", err)
	}

	// Send the request and wait for the response
	subject := "user.delete.service"
	resp, err := nc.Request(subject, reqDataBytes, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// Parse the response
	var respData deleteUserResponse
	if err := json.Unmarshal(resp.Data, &respData); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	// Check if there was an error in the response
	if respData.Error != "" {
		fmt.Printf("Error: %s\n", respData.Error)
	} else {
		fmt.Printf("Success: %s\n", respData.Message)
	}
}
