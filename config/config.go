package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	AppName      string
	AppPort      int
	LogLevel     string
	Environment  string
	DatabasePath string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	// Set default values
	config := &Config{
		AppName:      getEnv("APP_NAME", "Todo API"),
		AppPort:      getEnvAsInt("APP_PORT", 3000),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		DatabasePath: getEnv("DATABASE_PATH", "data/todo.db"),
	}

	return config, nil
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
