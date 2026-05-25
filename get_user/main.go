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

func getUserFromDB(db *sql.DB, userID int) (*User, error) {
	var user User
	query := `SELECT id, username FROM users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, err
	}
	return &user, nil
}

func handleGetUser(req micro.Request) {
	var requestData map[string]int

	if err := json.Unmarshal(req.Data(), &requestData); err != nil {
		log.Printf("Failed to parse request: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Invalid request"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	id, ok := requestData["id"]
	if !ok {
		if err := req.RespondJSON(map[string]string{"error": "User ID not provided"}); err != nil {
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

	user, err := getUserFromDB(db, id)
	if err != nil {
		log.Printf("Failed to retrieve user: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "User not found"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	config := micro.Config{
		Name:    "get_user-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "users.get",
			Handler: micro.HandlerFunc(handleGetUser),
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
