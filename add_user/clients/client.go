package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	requestData := `{"username": "testuser"}`
	response, err := nc.Request("users.add", []byte(requestData), 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	var responseData map[string]string
	if err := json.Unmarshal(response.Data, &responseData); err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %v", responseData)
}
