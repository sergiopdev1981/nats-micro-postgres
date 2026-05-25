package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestSendUpdateRequest(t *testing.T) {
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		t.Fatalf("Failed to connect to NATS server: %v", err)
	}
	defer nc.Close()

	if _, err := nc.Subscribe("users.update", func(m *nats.Msg) {
		response := map[string]string{
			"message": "User updated successfully",
		}
		responseData, _ := json.Marshal(response)
		if err := m.Respond(responseData); err != nil {
			log.Printf("Failed to respond: %v", err)
		}
	}); err != nil {
		log.Printf("Error subscribing to subject: %v", err)
	}

	sendRequest := func() (map[string]string, error) {
		requestData := map[string]interface{}{
			"id":       1,
			"username": "updateduser",
		}
		data, _ := json.Marshal(requestData)
		response, err := nc.Request("users.update", data, 2*time.Second)
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

	if responseData["message"] != "User updated successfully" {
		t.Errorf("Expected success message, got %v", responseData["message"])
	}
}
