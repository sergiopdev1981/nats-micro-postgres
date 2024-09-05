package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func sendRequest() {
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Send request and wait for a response
	requestData := `{"username": "testuser"}`
	response, err := nc.Request("users.add.service", []byte(requestData), 10*time.Second) // 10-second timeout
	if err != nil {
		log.Fatal(err)
	}

	// Handle response
	var responseData map[string]string
	if err := json.Unmarshal(response.Data, &responseData); err != nil {
		log.Fatal(err)
	}

	log.Printf("Received response for user%v: %v", responseData["user_id"], responseData["message"])
}

func main() {
	sendRequest()
}
