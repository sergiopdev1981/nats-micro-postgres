package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func sendRequest() {
	// Connect to the running NATS server
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Send request to get all users and wait for a response
	response, err := nc.Request("users.get.service", nil, 10*time.Second) // No payload needed
	if err != nil {
		log.Fatal(err)
	}

	// Handle response
	var users []User
	if err := json.Unmarshal(response.Data, &users); err != nil {
		log.Fatal(err)
	}

	// Log the received users
	log.Printf("Received %d users:", len(users))
	for _, user := range users {
		log.Printf("User ID: %d, Username: %s", user.ID, user.Username)
	}
}

func main() {
	sendRequest()
}
