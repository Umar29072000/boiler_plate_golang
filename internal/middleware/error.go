package middleware

import (
	"boiler_plate_be_golang/pkg/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles errors globally
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			// Log error
			log.Printf("Error: %v", err)

			// Handle Fiber errors
			if e, ok := err.(*fiber.Error); ok {
				return utils.ErrorResponse(c, e.Code, e.Message, nil)
			}

			// Default to 500 Internal Server Error
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
		}

		return nil
	}
}
