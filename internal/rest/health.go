package rest

import (
	"net/http"

	"boiler_plate_be_golang/internal/rest/response"

	"github.com/gofiber/fiber/v2"
)

// healthHandler handles health check requests
type healthHandler struct{}

// InitHealthHandler initializes health check routes
func InitHealthHandler(e fiber.Router) {
	handler := &healthHandler{}

	healthGroup := e.Group("/health")
	healthGroup.Get("", handler.HealthCheck)
	healthGroup.Get("/ready", handler.ReadinessCheck)
	healthGroup.Get("/live", handler.LivenessCheck)
}

// HealthCheck handles basic health check
func (h *healthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"status": "healthy",
		},
	})
}

// ReadinessCheck handles readiness probe
func (h *healthHandler) ReadinessCheck(c *fiber.Ctx) error {
	// TODO: Add actual readiness checks (database, redis, etc.)
	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "READY",
		Data: map[string]interface{}{
			"status": "ready",
		},
	})
}

// LivenessCheck handles liveness probe
func (h *healthHandler) LivenessCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "ALIVE",
		Data: map[string]interface{}{
			"status": "alive",
		},
	})
}
