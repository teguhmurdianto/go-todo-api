# Go Todo API

A high-performance Todo API built with Go, Fiber, and SQLite following cloud-native production-grade standards.

## Features

- RESTful API for managing todo items
- Clean architecture with separation of concerns
- SQLite database for data persistence
- Swagger documentation
- Middleware for security, logging, and error handling
- Graceful shutdown
- Environment-based configuration

## Project Structure

```
.
├── cmd
│   └── api             # Application entry point
├── config              # Configuration management
├── internal
│   ├── app
│   │   ├── handlers    # HTTP handlers
│   │   ├── models      # Data models
│   │   ├── repositories # Data access layer
│   │   └── services    # Business logic
│   ├── database        # Database connection and migrations
│   └── middleware      # HTTP middleware
├── pkg
│   └── utils           # Utility functions
└── scripts             # Helper scripts
```

## Prerequisites

- Go 1.16+
- SQLite3

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-todo-api.git
cd go-todo-api
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file in the root directory (optional):

```
APP_NAME=Todo API
APP_PORT=3000
LOG_LEVEL=info
ENVIRONMENT=development
DATABASE_PATH=data/todo.db
```

## Setup and Running the API

### Database Migration

Before running the API, you need to set up the database by running the migration script:

```bash
# Run the migration script directly
go run scripts/migrate.go

# Or use the Makefile command
make migrate
```

This will create the SQLite database file and set up the required tables.

### Generate Swagger Documentation

To generate the Swagger documentation for the API:

```bash
# Install the Swagger tool (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate the Swagger docs
swag init -g cmd/api/main.go

# Or use the Makefile command
make swagger
```

### Running the API

After setting up the database and generating the Swagger documentation, you can run the API:

```bash
# Run directly with Go
go run cmd/api/main.go

# Or build and run
go build -o todo-api cmd/api/main.go
./todo-api

# Or use the Makefile command
make run
```

### Development Mode with Hot Reload

For development, you can use the air tool for hot reloading:

```bash
# Install air (if not already installed)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air

# Or use the Makefile command
make dev
```

## API Endpoints

| Method | Endpoint      | Description                                |
|--------|---------------|--------------------------------------------|
| GET    | /health       | Health check endpoint                      |
| GET    | /swagger/*    | Swagger documentation                      |
| POST   | /api/v1/todos | Create a new todo                          |
| GET    | /api/v1/todos | Get all todos (optional ?completed=true/false) |
| GET    | /api/v1/todos/:id | Get a specific todo by ID                 |
| PATCH  | /api/v1/todos/:id | Update a todo                             |
| DELETE | /api/v1/todos/:id | Delete a todo                             |

## API Requests and Responses

### Create Todo

**Request:**

```json
POST /api/v1/todos
{
  "title": "Complete project",
  "description": "Finish the Go Todo API project",
  "priority": 2,
  "due_date": "2023-12-31T23:59:59Z"
}
```

**Response:**

```json
Status: 201 Created
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete project",
  "description": "Finish the Go Todo API project",
  "completed": false,
  "priority": 2,
  "due_date": "2023-12-31T23:59:59Z",
  "created_at": "2023-04-01T12:00:00Z",
  "updated_at": "2023-04-01T12:00:00Z"
}
```

### Get All Todos

**Request:**

```
GET /api/v1/todos
```

**Response:**

```json
Status: 200 OK
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete project",
    "description": "Finish the Go Todo API project",
    "completed": false,
    "priority": 2,
    "due_date": "2023-12-31T23:59:59Z",
    "created_at": "2023-04-01T12:00:00Z",
    "updated_at": "2023-04-01T12:00:00Z"
  }
]
```

### Update Todo

**Request:**

```json
PATCH /api/v1/todos/550e8400-e29b-41d4-a716-446655440000
{
  "completed": true
}
```

**Response:**

```json
Status: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete project",
  "description": "Finish the Go Todo API project",
  "completed": true,
  "priority": 2,
  "due_date": "2023-12-31T23:59:59Z",
  "created_at": "2023-04-01T12:00:00Z",
  "updated_at": "2023-04-01T12:05:00Z"
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Generate and Access Swagger Documentation

Install swag:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate docs:

```bash
swag init -g cmd/api/main.go
```

After generating the documentation and starting the server, you can access the Swagger UI at:

```
http://localhost:3000/swagger/
```

This provides an interactive UI where you can explore and test all API endpoints.

## Deployment

### Docker

A Dockerfile is provided to containerize the application:

```bash
# Build the Docker image
docker build -t go-todo-api .

# Run the container
docker run -p 3000:3000 go-todo-api

# Or use Docker Compose
docker-compose up -d
```

### Kubernetes

Kubernetes manifests are available in the `k8s` directory for deploying to a Kubernetes cluster:

```bash
# Apply the Kubernetes manifests
kubectl apply -f k8s/deployment.yaml
```

## Project Commands (Makefile)

The project includes a Makefile with various commands to simplify development:

```bash
# Build the application
make build

# Run the application
make run

# Run in development mode with hot reload
make dev

# Run database migrations
make migrate

# Run tests
make test

# Generate Swagger documentation
make swagger

# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Show available commands
make help
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
