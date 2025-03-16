package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/teguh/go-todo-api/config"
	"github.com/teguh/go-todo-api/docs"
	"github.com/teguh/go-todo-api/internal/app/handlers"
	"github.com/teguh/go-todo-api/internal/database"
	"github.com/teguh/go-todo-api/internal/middleware"
)

// @title Todo API
// @version 1.0
// @description A high-performance Todo API built with Go and Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
func main() {
	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "A high-performance Todo API built with Go and Fiber"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := database.Initialize(cfg.DatabasePath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: customErrorHandler,
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// API routes
	api := app.Group("/api/v1")

	// Register handlers
	todoHandler := handlers.NewTodoHandler()
	todoHandler.RegisterRoutes(api)

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Root route - redirect to Swagger
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/", fiber.StatusMovedPermanently)
	})

	// Handle graceful shutdown
	go handleShutdown(app)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.AppPort)
	log.Printf("Starting %s server on %s in %s mode", cfg.AppName, addr, cfg.Environment)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// customErrorHandler handles errors thrown by Fiber
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default error response
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		// Log non-Fiber errors
		log.Printf("Error: %v", err)
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

// handleShutdown handles graceful shutdown
func handleShutdown(app *fiber.App) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
}
