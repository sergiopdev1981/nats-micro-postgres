package main

import (
	"database/sql"
	"encoding/json"
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

func getAllUsersFromDB(db *sql.DB) ([]User, error) {
	var users []User
	query := `SELECT id, username FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func handleGetUsers(req micro.Request) {
	db, err := connectDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database connection error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	defer db.Close()

	users, err := getAllUsersFromDB(db)
	if err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Database query error"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	responseData, err := json.Marshal(users)
	if err != nil {
		log.Printf("Failed to marshal users: %v", err)
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
		Name:    "get_users-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "users.list",
			Handler: micro.HandlerFunc(handleGetUsers),
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
