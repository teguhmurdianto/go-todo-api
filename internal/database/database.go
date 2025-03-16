package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
var DB *sql.DB

// Initialize sets up the database connection and creates tables if they don't exist
func Initialize(dbPath string) error {
	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Store the database connection
	DB = db

	// Create tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// createTables creates the necessary tables if they don't exist
func createTables() error {
	// Create todos table
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN NOT NULL DEFAULT 0,
		priority INTEGER NOT NULL DEFAULT 0,
		due_date TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	return nil
}
