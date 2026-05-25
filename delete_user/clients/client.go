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

type deleteUserRequest struct {
	ID int `json:"id"`
}

type deleteUserResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	userID := 1

	reqData := deleteUserRequest{ID: userID}
	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		log.Fatalf("Failed to marshal request data: %v", err)
	}

	resp, err := nc.Request("users.delete", reqDataBytes, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	var respData deleteUserResponse
	if err := json.Unmarshal(resp.Data, &respData); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if respData.Error != "" {
		fmt.Printf("Error: %s\n", respData.Error)
	} else {
		fmt.Printf("Success: %s\n", respData.Message)
	}
}
