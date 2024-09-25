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
	ID int `json:"id"`
}

// connecting to database
func connectDB() (*sql.DB, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable" + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	return sql.Open("postgres", connStr)
}

// delete user, this is executed when a request is called
func deleteUserDB(db *sql.DB, userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	// check if row affected
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
	// this is where I save users/service data
	var requestData map[string]int
	/*
		"req.Data()" has all user request data
		"&requestData" is a pointer to the struct and here is where the data (unmarshalled) will
		be stored
		"req.RespondJSON(map[string]string{"error": "invalid request"})"
		  "req.RespondJSON" sends a response to the configured subject/s
		  "map[string]string" creates a container to save a key-value object
		  "{"error": "invalid request"}" the key-value object

	*/
	if err := json.Unmarshal(req.Data(), &requestData); err != nil {
		log.Printf("Failed to parse request %v", err)
		req.RespondJSON(map[string]string{"error": "invalid request"})
		return
	}

	// here is the user data extraction
	/*
		"userID, err := requestData["id"]"
		  "userID" the var where the user requested data will be stored
		  "ok" same as userID but only to save if ok ocurred (only if function returns it)
		  "requestData["id"]" here is where unmarshalled data was stored
		    ["id"] this is a value that has to be on request data
	*/
	userID, ok := requestData["id"]
	if !ok {
		req.RespondJSON(map[string]string{"error": "User ID is not provided"})
		return
	}

	// connect to database
	db, err := connectDB()
	if err != nil {
		log.Printf("Failed to connect to database %v", err)
		req.RespondJSON(map[string]string{"error": "Database connection error"})
		return
	}
	// with this db will close eventually
	defer db.Close()

	// here is where the created function is used
	/*
		  "err = deleteUserDB(db, userID)"
		    "err" saves an error (only if function returns it)
			"deleteUserDB(db, userID)" function call passing the db connfiguration and a param
	*/
	err = deleteUserDB(db, userID)
	if err != nil {
		log.Printf("Failed to delete user from database %v", err)
		req.RespondJSON(map[string]string{"error": err.Error()})
		return
	}

	// responds ok to client if magic happened
	req.RespondJSON(map[string]string{"message": "User deleted succesfully"})
}

func main() {
	// loading env vars
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file:  %v", err)
	}

	// creating a nats connection
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// adding nats-micro service config
	deleteUserConfig := micro.Config{
		Name:    "delete_user-service",
		Version: "0.1.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "user.delete.service",
			Handler: micro.HandlerFunc(handleDeleteUser), // Handle Delete User Request
		},
	}

	// adding nats-micro service
	deleteUserSvc, err := micro.AddService(nc, deleteUserConfig)
	if err != nil {
		log.Fatal(err)
	}

	// adding a defer function to handle deleteUserSvc function end
	// if it's not stoppping then send message
	defer func() {
		if err := deleteUserSvc.Stop(); err != nil {
			log.Printf("Failed to stop delete user service: %v", err)
		}
	}()

	select {}
}
