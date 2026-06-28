package controllers

import (
	"boiler_plate_be_golang/pkg/redis"
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MetricsController handles metrics requests
type MetricsController struct{}

// NewMetricsController creates new metrics controller
func NewMetricsController() *MetricsController {
	return &MetricsController{}
}

// GetMetrics returns application runtime metrics
func (m *MetricsController) GetMetrics(c *fiber.Ctx) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return c.JSON(fiber.Map{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"runtime": fiber.Map{
			"go_version":   runtime.Version(),
			"goroutines":   runtime.NumGoroutine(),
			"cpu_count":    runtime.NumCPU(),
			"memory_alloc": bytesToMB(memStats.Alloc),
			"memory_total": bytesToMB(memStats.TotalAlloc),
			"memory_sys":   bytesToMB(memStats.Sys),
			"gc_count":     memStats.NumGC,
			"last_gc":      formatUnixNano(memStats.LastGC),
		},
		"services": fiber.Map{
			"redis_connected": redis.IsConnected(),
		},
	})
}

func bytesToMB(bytes uint64) string {
	return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
}

func formatUnixNano(value uint64) string {
	if value == 0 {
		return "never"
	}
	return time.Unix(0, int64(value)).UTC().Format(time.RFC3339)
}
