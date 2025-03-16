package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/teguh/go-todo-api/internal/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Get database path from environment or use default
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "data/todo.db"
	}

	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Initialize database
	log.Println("Running database migrations...")
	if err := database.Initialize(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Migrations completed successfully!")
}
