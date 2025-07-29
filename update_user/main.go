package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the pq driver
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

// create a struct to handlle client's data
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// connecting to database
func connectDB() (*sql.DB, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable" + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	return sql.Open("postgres", connStr)
}

// add user function, this is executed when a request is called
func addUserToDB(db *sql.DB, username string) (int, error) {
	var userID int
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id`
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	// here I return the database return
	return userID, nil
}

// this is the main function. here I get and parse client's data, send the data to the adduser function and send the response to the client
func handleAddUser(req micro.Request) {
	// create a "container" for the client's data
	var requestData map[string]string
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
	username, ok := requestData["username"]
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
			log.Printf("Failed to map error connecting to database reponse: %v", err)
		}
		return
	}
	defer db.Close()

	// Add user to database
	userID, err := addUserToDB(db, username)
	if err != nil {
		// if can't add user then send message
		log.Printf("Failed to add user to database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database insertion error"}); err != nil {
			log.Printf("Failed to response insertion data error: %v", err)
		}
		return
	}

	// Prepare and send response if user add worked ok
	response := map[string]string{
		"message": "User successfully added!!!",
		"user_id": fmt.Sprintf("%d", userID), // Ensure userID is converted to string
	}

	// Debug log the response being sent
	log.Printf("Sending response for user%v: %v", userID, response["message"])

	// Send response
	if err := req.RespondJSON(response); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func main() {
	// creates a nats conenction
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// loads .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// nats microservice config
	config := micro.Config{
		// service name
		Name: "add_user-service",
		// service version
		Version: "0.1.0",
		// micro service endpoint config
		Endpoint: &micro.EndpointConfig{
			// subject where magic happens :)
			Subject: "users.add.service",
			// function which make the magic happens :)
			Handler: micro.HandlerFunc(handleAddUser),
		},
	}
	// adding service passing nats connection and service config
	svc, err := micro.AddService(nc, config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := svc.Stop(); err != nil {
			log.Printf("Failed to stop service: %v", err)
		}
	}()

	select {} // Keeps the service running
}
