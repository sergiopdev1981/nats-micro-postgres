package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestSendRequest(t *testing.T) {
	// Connect to the running NATS server
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		t.Fatalf("Failed to connect to NATS server: %v", err)
	}
	defer nc.Close()

	// Subscribe to the request subject and provide a mock response
	if _, err := nc.Subscribe("users.add.service", func(m *nats.Msg) {
		response := map[string]string{
			"message": "User successfully added",
			"user_id": "123",
		}
		responseData, _ := json.Marshal(response)
		if err := m.Respond(responseData); err != nil {
			log.Printf("Failed to respond data: %v", err)
		}
	}); err != nil {
		log.Printf("Error trying to subscribe to subject: %v", err)
	}

	// Function to test the request
	sendRequest := func() (map[string]string, error) {
		response, err := nc.Request("users.add.service", []byte(`{"username": "testuser"}`), 2*time.Second)
		if err != nil {
			return nil, err
		}

		var responseData map[string]string
		if err := json.Unmarshal(response.Data, &responseData); err != nil {
			return nil, err
		}

		return responseData, nil
	}

	responseData, err := sendRequest()
	if err != nil {
		t.Fatalf("Failed to get response: %v", err)
	}

	expectedResponse := map[string]string{
		"message": "User successfully added",
		"user_id": "123",
	}

	if responseData["message"] != expectedResponse["message"] {
		t.Errorf("Expected message %v, got %v", expectedResponse["message"], responseData["message"])
	}
	if responseData["user_id"] != expectedResponse["user_id"] {
		t.Errorf("Expected user_id %v, got %v", expectedResponse["user_id"], responseData["user_id"])
	}
}
