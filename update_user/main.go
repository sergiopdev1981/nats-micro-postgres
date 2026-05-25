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

func updateUserInDB(db *sql.DB, userID int, username string) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	result, err := db.Exec(query, username, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func handleUpdateUser(req micro.Request) {
	var requestData map[string]interface{}

	if err := json.Unmarshal(req.Data(), &requestData); err != nil {
		log.Printf("Failed to parse request: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "Invalid request"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	idFloat, ok := requestData["id"].(float64)
	if !ok {
		if err := req.RespondJSON(map[string]string{"error": "User ID not provided or invalid"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	userID := int(idFloat)

	username, ok := requestData["username"].(string)
	if !ok || username == "" {
		if err := req.RespondJSON(map[string]string{"error": "Username not provided or invalid"}); err != nil {
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

	if err := updateUserInDB(db, userID, username); err != nil {
		log.Printf("Failed to update user: %v", err)
		if err := req.RespondJSON(map[string]string{"error": "User not found or update failed"}); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}

	if err := req.RespondJSON(map[string]string{"message": "User updated successfully"}); err != nil {
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
		Name:    "update_user-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "users.update",
			Handler: micro.HandlerFunc(handleUpdateUser),
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
