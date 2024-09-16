package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// connecting to database
func connectDB() (*sql.DB, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable" + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	return sql.Open("postgres", connStr)
}

func getUserFromDB(db *sql.DB, userID int) (*User, error) {
	var user User
	query := `SELECT id, username FROM users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// if no user is found, return nil and a custom error
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, err
	}
	return &user, nil
}

// Function to handle the request to get a single user by ID
func handleGetUser(req micro.Request) {
	// create a "container" for the client's data
	var requestData map[string]int
	// json.Unmarshal take data sent by client and converts it into go data, in this case req.Data() has client's data
	if err := json.Unmarshal(req.Data(), &requestData); err != nil {
		// if can't logs and send a response to client
		log.Printf("Failed to parse request: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Invalid request"}); err != nil {
			log.Printf("Failed to map request data: %v", err)
		}
		return
	}

	// Extract username from request -> from the "container" created
	id, ok := requestData["id"]
	if !ok {
		// if not data founded send message
		if err := req.RespondJSON(map[string]string{"error": "Username not provided"}); err != nil {
			log.Printf("Failed to map username not provided error response: %v", err)
		}
		return
	}

	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database connection error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	defer db.Close()

	// Retrieve the user from the database
	user, err := getUserFromDB(db, id)
	if err != nil {
		log.Printf("Failed to retrieve user from database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "User not found"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	// Convert the user data to JSON and send the response
	responseData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to marshal user: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Internal server error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	if err := req.Respond(responseData); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func main() {
	// creates a nats connection
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// loads .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// nats microservice config for get user
	userConfig := micro.Config{
		Name:    "get_user-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "user.get.service",
			Handler: micro.HandlerFunc(handleGetUser),
		},
	}
	// adding service passing nats connection and service config
	userSvc, err := micro.AddService(nc, userConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := userSvc.Stop(); err != nil {
			log.Printf("Failed to stop user service: %v", err)
		}
	}()
	select {} // Keeps the service running
}
