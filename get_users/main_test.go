package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestHandleGetUser(t *testing.T) {
	// Connect to the running NATS server
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		t.Fatalf("Failed to connect to NATS server: %v", err)
	}
	defer nc.Close()

	// Mocking the database response
	mockUser := User{ID: 1, Username: "mockuser"}
	mockDB := func(userID int) (*User, error) {
		if userID == mockUser.ID {
			return &mockUser, nil
		}
		return nil, sql.ErrNoRows
	}

	// Subscribe to the request subject and provide a mock response
	if _, err := nc.Subscribe("user.get.service", func(m *nats.Msg) {
		// Unmarshal the request data
		var requestData map[string]int
		if err := json.Unmarshal(m.Data, &requestData); err != nil {
			log.Printf("Failed to parse request: %v", err)
			return
		}

		// Get user ID from the request
		userID, ok := requestData["id"]
		if !ok {
			// Respond with an error if ID is missing
			m.Respond([]byte(`{"error":"ID not provided"}`))
			return
		}

		// Call the mockDB function to get the user
		user, err := mockDB(userID)
		if err != nil {
			// Respond with an error if the user is not found
			m.Respond([]byte(`{"error":"User not found"}`))
			return
		}

		// Marshal the user data to JSON and respond
		responseData, _ := json.Marshal(user)
		m.Respond(responseData)
	}); err != nil {
		log.Fatalf("Failed to subscribe to subject: %v", err)
	}

	// Function to send the request and receive the response
	sendRequest := func(userID int) (*User, error) {
		// Prepare request data
		requestData := map[string]int{"id": userID}
		data, _ := json.Marshal(requestData)

		// Send the request to the NATS server
		response, err := nc.Request("user.get.service", data, 2*time.Second)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data
		var user User
		if err := json.Unmarshal(response.Data, &user); err != nil {
			return nil, err
		}

		return &user, nil
	}

	// Test the function with a valid user ID
	t.Run("ValidUserID", func(t *testing.T) {
		user, err := sendRequest(mockUser.ID)
		if err != nil {
			t.Fatalf("Failed to get response: %v", err)
		}

		// Check if the returned user matches the expected mock user
		if user.ID != mockUser.ID || user.Username != mockUser.Username {
			t.Errorf("Expected user %+v, got %+v", mockUser, user)
		}
	})

	// Test the function with an invalid user ID
	t.Run("InvalidUserID", func(t *testing.T) {
		_, err := sendRequest(999) // Invalid user ID
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
	})
}
