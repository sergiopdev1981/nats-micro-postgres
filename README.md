# Microservices Project

[![Go Report Card](https://goreportcard.com/badge/github.com/sergiopdev1981/labs)](https://goreportcard.com/report/github.com/sergiopdev1981/labs)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This repository contains a collection of microservices written in Go. Each microservice is designed to handle a specific functionality and can be deployed independently.

## Microservices Overview

### 1. `add_user`
- **Description**: Handles the addition of new users to the system.
- **Features**:
  - Adds user data to the database.
  - Validates input data.
- **Dependencies**:
  - `github.com/joho/godotenv`
  - `github.com/nats-io/nats.go`

### 2. `delete_user`
- **Description**: Manages the deletion of users from the system.
- **Features**:
  - Deletes user data from the database.
  - Enhanced error handling for JSON responses.
- **Dependencies**:
  - `github.com/DATA-DOG/go-sqlmock`
  - `github.com/joho/godotenv`
  - `github.com/nats-io/nats.go`

### 3. `get_user`
- **Description**: Retrieves information about a specific user.
- **Features**:
  - Fetches user data by ID.
  - Returns data in JSON format.
- **Dependencies**:
  - `github.com/joho/godotenv`
  - `github.com/nats-io/nats.go`

### 4. `get_users`
- **Description**: Retrieves a list of all users.
- **Features**:
  - Fetches all user data from the database.
  - Supports pagination.
- **Dependencies**:
  - `github.com/joho/godotenv`
  - `github.com/nats-io/nats.go`

### 5. `update_user`
- **Description**: Handles updates to user information.
- **Features**:
  - Updates user data in the database.
  - Validates input data.
- **Dependencies**:
  - `github.com/joho/godotenv`
  - `github.com/nats-io/nats.go`

## Prerequisites

To run these microservices, ensure you have the following services running:

### PostgreSQL
- **Environment Variables**:
  - `POSTGRES_USER`: Database username (e.g., `user1`)
  - `POSTGRES_PASSWORD`: Database password (e.g., `password1`)
  - `POSTGRES_DB`: Database name (e.g., `db_service1`)
  - `POSTGRES_HOST`: Hostname or IP address (e.g., `localhost`)
  - `POSTGRES_PORT`: Port number (e.g., `5431`)

### NATS
- **Environment Variables**:
  - `NATS_URL`: NATS server URL (e.g., `localhost:4222`)

Ensure these variables are set in the `.env` files for each microservice.

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/sergiopdev1981/labs.git
   cd labs/microservices
   ```

2. Run the linting and dependency setup script:
   ```bash
   ./run_lint.sh
   ```

3. Start each microservice individually by navigating to its directory and running:
   ```bash
   go run main.go
   ```

