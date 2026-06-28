package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// InitHealthHandler initializes health check routes
func InitHealthHandler(e fiber.Router) {
	handler := &HealthHandler{}

	healthGroup := e.Group("/health")
	healthGroup.Get("", handler.HealthCheck)
	healthGroup.Get("/ready", handler.ReadinessCheck)
	healthGroup.Get("/live", handler.LivenessCheck)
}

// HealthCheck handles basic health check
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "OK",
		"data": fiber.Map{
			"status": "healthy",
		},
	})
}

// ReadinessCheck handles readiness probe
func (h *HealthHandler) ReadinessCheck(c *fiber.Ctx) error {
	// TODO: Add actual readiness checks (database, redis, etc.)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "READY",
		"data": fiber.Map{
			"status": "ready",
		},
	})
}

// LivenessCheck handles liveness probe
func (h *HealthHandler) LivenessCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "ALIVE",
		"data": fiber.Map{
			"status": "alive",
		},
	})
}
