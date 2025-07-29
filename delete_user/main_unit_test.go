package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteUserDB(t *testing.T) {
	// Create mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Define the expected query and result
	userID := 1
	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	// Call the function
	err = deleteUserDB(db, userID)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestDeleteUserDB_NoRowsAffected(t *testing.T) {
	// Create mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Define the expected query and result
	userID := 1
	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected

	// Call the function
	err = deleteUserDB(db, userID)
	if err == nil {
		t.Fatalf("Expected error for no rows affected, but got nil")
	}

	expectedError := "No user found with ID 1"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
	}
}
