package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SendError sends an error response
func SendError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// SendSuccess sends a success response
func SendSuccess(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}
