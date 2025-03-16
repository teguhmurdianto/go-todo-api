package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
)

// SetupMiddleware sets up all middleware for the application
func SetupMiddleware(app *fiber.App) {
	// Recover from panics
	app.Use(recover.New())

	// Security headers
	app.Use(helmet.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Logger
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Request ID
	app.Use(func(c *fiber.Ctx) error {
		// Add request ID if not present
		if c.Get("X-Request-ID") == "" {
			c.Set("X-Request-ID", time.Now().Format("20060102150405")+"-"+c.IP())
		}
		return c.Next()
	})
}
