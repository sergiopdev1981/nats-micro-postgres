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
	ID int `json:"id"`
}

func connectDB() (*sql.DB, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable" + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	return sql.Open("postgres", connStr)
}

func deleteUserDB(db *sql.DB, userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No user found with ID %d", userID)
	}

	return nil
}

func handleDeleteUser(req micro.Request) {
	var requestData map[string]int

	if err := json.Unmarshal(req.Data(), &requestData); err != nil {
		log.Printf("Failed to parse request: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Invalid request"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	userID, ok := requestData["id"]
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

	err = deleteUserDB(db, userID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		if err := req.RespondJSON(map[string]string{"error": err.Error()}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	if err := req.RespondJSON(map[string]string{"message": "User deleted successfully"}); err != nil {
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
		Name:    "delete_user-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "users.delete",
			Handler: micro.HandlerFunc(handleDeleteUser),
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
