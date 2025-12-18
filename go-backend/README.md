# User API - Go Backend

A RESTful API built with Go for managing users with dynamic age calculation.

## Tech Stack

- **Framework**: GoFiber
- **Database**: PostgreSQL + SQLC
- **Logging**: Uber Zap
- **Validation**: go-playground/validator

## Project Structure

```
/cmd/server/main.go          # Application entry point
/config/                      # Configuration management
/db/migrations/               # SQL migration files
/db/sqlc/                     # SQLC queries and generated code
/internal/
├── handler/                  # HTTP handlers
├── repository/               # Data access layer
├── service/                  # Business logic
├── routes/                   # Route definitions
├── middleware/               # Custom middleware
├── models/                   # Data models
└── logger/                   # Logging setup
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Running with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app
```

### Running Locally

1. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

2. **Start PostgreSQL**
   ```bash
   # Using Docker
   docker run -d \
     --name postgres \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=userdb \
     -p 5432:5432 \
     postgres:15-alpine
   ```

3. **Run migrations**
   ```bash
   psql -h localhost -U postgres -d userdb -f db/migrations/001_init.sql
   ```

4. **Install dependencies & run**
   ```bash
   go mod download
   go run cmd/server/main.go
   ```

### Generate SQLC Code

```bash
# Install SQLC
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code
cd db/sqlc && sqlc generate
```

## API Endpoints

| Method | Endpoint        | Description      |
|--------|-----------------|------------------|
| POST   | /api/v1/users   | Create user      |
| GET    | /api/v1/users   | List all users   |
| GET    | /api/v1/users/:id | Get user by ID |
| PUT    | /api/v1/users/:id | Update user    |
| DELETE | /api/v1/users/:id | Delete user    |
| GET    | /health         | Health check     |

## API Examples

### Create User
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "dob": "1990-05-10"}'
```

Response:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User (with calculated age)
```bash
curl http://localhost:3000/api/v1/users/1
```

Response:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 34
}
```

### List Users (with pagination)
```bash
curl "http://localhost:3000/api/v1/users?page=1&page_size=10"
```

### Update User
```bash
curl -X PUT http://localhost:3000/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated", "dob": "1991-03-15"}'
```

### Delete User
```bash
curl -X DELETE http://localhost:3000/api/v1/users/1
```

## Running Tests

```bash
go test ./... -v
```

## Environment Variables

| Variable    | Description          | Default    |
|-------------|----------------------|------------|
| PORT        | Server port          | 3000       |
| LOG_LEVEL   | Log level            | info       |
| DB_HOST     | Database host        | localhost  |
| DB_PORT     | Database port        | 5432       |
| DB_USER     | Database user        | postgres   |
| DB_PASSWORD | Database password    | postgres   |
| DB_NAME     | Database name        | userdb     |
| DB_SSLMODE  | SSL mode             | disable    |

## Features

- ✅ CRUD operations for users
- ✅ Dynamic age calculation from DOB
- ✅ Input validation with go-playground/validator
- ✅ Structured logging with Uber Zap
- ✅ Request ID middleware for tracing
- ✅ Pagination support
- ✅ Docker support
- ✅ Unit tests for age calculation
