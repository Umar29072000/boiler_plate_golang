package controllers

import (
	"boiler_plate_be_golang/internal/config"
	"boiler_plate_be_golang/internal/database"
	"boiler_plate_be_golang/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2"
)

var serverStartTime = time.Now()

// HealthController handles health check requests
type HealthController struct{}

// NewHealthController creates new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// BasicHealth handles basic health check
func (h *HealthController) BasicHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Server is running",
	})
}

// DetailedHealth handles detailed health check with service status
func (h *HealthController) DetailedHealth(c *fiber.Ctx) error {
	uptime := time.Since(serverStartTime)
	
	// Check database connection
	dbStatus := "healthy"
	dbError := ""
	if db, err := database.DB.DB(); err != nil {
		dbStatus = "unhealthy"
		dbError = err.Error()
	} else if err := db.Ping(); err != nil {
		dbStatus = "unhealthy"
		dbError = err.Error()
	}

	// Check Redis connection
	redisStatus := "healthy"
	redisError := ""
	if !redis.IsConnected() {
		redisStatus = "unhealthy"
		redisError = "Redis not connected"
	}

	// Determine overall status
	overallStatus := "OK"
	statusCode := fiber.StatusOK
	
	if dbStatus == "unhealthy" {
		overallStatus = "ERROR"
		statusCode = fiber.StatusServiceUnavailable
	} else if redisStatus == "unhealthy" {
		overallStatus = "DEGRADED"
		// Redis is optional, so we stay at 200
	}

	response := fiber.Map{
		"status":    overallStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    uptime.String(),
		"services": fiber.Map{
			"database": fiber.Map{
				"status": dbStatus,
				"type":   "PostgreSQL",
			},
			"redis": fiber.Map{
				"status": redisStatus,
				"type":   "Redis",
			},
		},
		"environment": config.App.App.Env,
		"version":     config.App.App.Name,
	}

	if dbError != "" {
		response["services"].(fiber.Map)["database"].(fiber.Map)["error"] = dbError
	}
	if redisError != "" {
		response["services"].(fiber.Map)["redis"].(fiber.Map)["error"] = redisError
	}

	return c.Status(statusCode).JSON(response)
}

// ReadinessProbe for Kubernetes readiness probe
func (h *HealthController) ReadinessProbe(c *fiber.Ctx) error {
	// Check if critical services are ready
	if db, err := database.DB.DB(); err != nil || db.Ping() != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"ready": false,
			"error": "Database not ready",
		})
	}

	return c.JSON(fiber.Map{
		"ready": true,
	})
}

// LivenessProbe for Kubernetes liveness probe
func (h *HealthController) LivenessProbe(c *fiber.Ctx) error {
	// Simple check that server is running
	return c.JSON(fiber.Map{
		"alive": true,
	})
}
