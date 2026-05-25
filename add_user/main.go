package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func connectDB() (*sql.DB, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable" + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	return sql.Open("postgres", connStr)
}

func addUserToDB(db *sql.DB, username string) (int, error) {
	var userID int
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id`
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func handleAddUser(req micro.Request) {
	var requestData map[string]string

	if err := json.Unmarshal(req.Data(), &requestData); err != nil {
		log.Printf("Failed to parse request: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Invalid request"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	username, ok := requestData["username"]
	if !ok {
		if err := req.RespondJSON(map[string]string{"error": "Username not provided"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	db, err := connectDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database connection error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	defer db.Close()

	userID, err := addUserToDB(db, username)
	if err != nil {
		log.Printf("Failed to add user to database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database insertion error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	response := map[string]string{
		"message": "User successfully added",
		"user_id": fmt.Sprintf("%d", userID),
	}

	if err := req.RespondJSON(response); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	config := micro.Config{
		Name:    "add_user-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "users.add",
			Handler: micro.HandlerFunc(handleAddUser),
		},
	}

	svc, err := micro.AddService(nc, config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := svc.Stop(); err != nil {
			log.Printf("Failed to stop service: %v", err)
		}
	}()

	select {}
}
