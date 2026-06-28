package middleware

import (
	"time"

	"boiler_plate_be_golang/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger creates a structured logging middleware
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)

		// Start timer
		start := time.Now()

		// Log request
		logger.Logger.Info().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("Request started")

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		logEvent := logger.Logger.Info().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Int("size", len(c.Response().Body()))

		// Add error if exists
		if err != nil {
			logEvent = logger.Logger.Error().
				Str("request_id", requestID).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Int("status", c.Response().StatusCode()).
				Dur("duration", duration).
				Err(err)
		}

		logEvent.Msg("Request completed")

		return err
	}
}

// RequestID middleware adds request ID to context
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context
		c.Locals("request_id", requestID)

		// Add to response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
