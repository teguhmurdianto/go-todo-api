# Build stage
FROM golang:1.24.1-alpine AS builder

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o todo-api ./cmd/api

# Build the migration tool
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o migrate ./scripts/migrate.go

# Final stage
FROM alpine:latest

# Install required packages
RUN apk --no-cache add ca-certificates tzdata sqlite

# Set working directory
WORKDIR /app

# Copy the binaries from builder
COPY --from=builder /app/todo-api .
COPY --from=builder /app/migrate .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 3000

# Run migration and start the application
CMD ["sh", "-c", "./migrate && ./todo-api"]
