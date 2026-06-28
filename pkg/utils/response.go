package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse sends success response
func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends error response
func ErrorResponse(c *fiber.Ctx, status int, message string, err interface{}) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}
