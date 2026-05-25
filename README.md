<p align="center">
  <img src="https://img.shields.io/badge/lang-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/messaging-NATS-27AAE1?style=for-the-badge&logo=nats&logoColor=white" alt="NATS">
  <img src="https://img.shields.io/badge/database-PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/architecture-Microservices-FF6F00?style=for-the-badge" alt="Microservices">
</p>

# NATS Microservices with PostgreSQL

A collection of microservices written in Go, communicating via NATS message broker and persisting data in PostgreSQL.

## Services

| Service | Subject | Description |
|---------|---------|-------------|
| **add_user** | `users.add` | Creates a new user with the given username |
| **get_user** | `users.get` | Retrieves a single user by ID |
| **get_users** | `users.list` | Lists all users |
| **update_user** | `users.update` | Updates a user's username by ID |
| **delete_user** | `users.delete` | Deletes a user by ID |

## Architecture

Each microservice runs independently and communicates exclusively through NATS subjects. Services do not expose HTTP endpoints — all requests are NATS request-reply messages.

```
Client ──> NATS ──> Microservice ──> PostgreSQL
```

Example request payload:
```json
{"username": "johndoe"}
```

## Prerequisites

### PostgreSQL
| Variable | Description | Default |
|----------|-------------|---------|
| `POSTGRES_USER` | Database username | `user1` |
| `POSTGRES_PASSWORD` | Database password | `password1` |
| `POSTGRES_DB` | Database name | `db_service1` |
| `POSTGRES_HOST` | Hostname | `localhost` |
| `POSTGRES_PORT` | Port | `5431` |

### NATS
| Variable | Description | Default |
|----------|-------------|---------|
| `NATS_URL` | NATS server URL | `localhost:4222` |

Copy the template config:
```bash
cp add_user/.env.example add_user/.env
# Repeat for each service or symlink a single .env
```

## How to Run

1. Make sure PostgreSQL and NATS are running.

2. Create the `users` table:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL
);
```

3. Start services (in separate terminals):
```bash
cd add_user && go run main.go
cd get_user && go run main.go
```

4. Send a request using one of the test clients:
```bash
cd add_user/clients && go run client.go
```

Or programmatically via NATS:
```go
response, err := nc.Request("users.add", []byte(`{"username":"alice"}`), 2*time.Second)
```

## Linting

```bash
./run_lint.sh
```
