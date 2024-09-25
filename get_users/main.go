package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the pq driver
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

// create here a struct to handle client data
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// connecting to database
func connectDB() (*sql.DB, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable" + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	return sql.Open("postgres", connStr)
}

// Function to retrieve all users from the database
func getAllUsersFromDB(db *sql.DB) ([]User, error) {
	// create an array to save response data after query database
	var users []User
	// query to database
	query := `SELECT id, username FROM users`
	// executing the query and saving the response to a variable
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// looping the rows after executing query
	for rows.Next() {
		// "instance" of the User struct to save data
		var user User
		// in every loop save data... using & to not copy the struct
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		// append the saved data in the array
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	// then return the array
	return users, nil
}

// Function to handle the request to get all users
func handleGetUsers(req micro.Request) {
	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		// if error then response the error to client, has to map as json
		if err := req.RespondJSON(map[string]string{"error": "Database connection error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	defer db.Close()

	// Retrieve all users from the database
	users, err := getAllUsersFromDB(db)
	if err != nil {
		log.Printf("Failed to retrieve users from database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database query error"}); err != nil {
			log.Printf("Failed to execute query: %v", err)
		}
		return
	}

	// Convert data from executed query to JSON and send the response
	responseData, err := json.Marshal(users)
	// if error send error to client
	if err != nil {
		log.Printf("Failed to marshal users: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Internal server error"}); err != nil {
			log.Printf("Failed to map users data: %v", err)
		}
		return
	}
	// if not error then send the json parsed data
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
	// nats microservice config
	config := micro.Config{
		// service name
		Name: "get_users-service",
		// service version
		Version: "0.1.0",
		// microservice endpoint config
		Endpoint: &micro.EndpointConfig{
			// subject where magic happens :)
			Subject: "users.get.service",
			// function which make the magic happens :)
			Handler: micro.HandlerFunc(handleGetUsers),
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
