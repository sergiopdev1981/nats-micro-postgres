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

func main() {
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	response, err := nc.Request("users.list", nil, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	if err := json.Unmarshal(response.Data, &users); err != nil {
		log.Fatal(err)
	}

	log.Printf("Received %d users:", len(users))
	for _, user := range users {
		log.Printf("User ID: %d, Username: %s", user.ID, user.Username)
	}
}
