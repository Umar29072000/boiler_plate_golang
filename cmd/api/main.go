package main

import (
	"boiler_plate_be_golang/internal/config"
	"boiler_plate_be_golang/internal/database"
	"boiler_plate_be_golang/internal/database/migrations"
	"boiler_plate_be_golang/internal/middleware"
	"boiler_plate_be_golang/internal/routes"
	"boiler_plate_be_golang/pkg/logger"
	"boiler_plate_be_golang/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		logger.Fatal("Failed to load config").Err(err).Send()
	}

	// Initialize logger (after config is loaded)
	logger.Init()

	logger.Info("Starting application").
		Str("app", config.App.App.Name).
		Str("env", config.App.App.Env).
		Send()

	// Connect to database
	if err := database.Connect(); err != nil {
		logger.Fatal("Failed to connect to database").Err(err).Send()
	}
	defer database.Close()

	// Connect to Redis (optional, non-fatal)
	if err := redis.Connect(); err != nil {
		logger.Warn("Redis connection failed, using fallback").Err(err).Send()
	}
	defer redis.Close()

	// Run migrations
	if err := migrations.Migrate(database.GetDB()); err != nil {
		logger.Fatal("Failed to run migrations").Err(err).Send()
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      config.App.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Compress())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.CORS())
	app.Use(middleware.DefaultRateLimiter())
	app.Use(middleware.XSSProtection())
	app.Use(middleware.InputSanitizer())
	app.Use(middleware.ErrorHandler())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := config.App.App.Port
	logger.Info("Server starting").
		Str("port", port).
		Str("environment", config.App.App.Env).
		Str("url", config.App.App.URL).
		Send()

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server").Err(err).Send()
	}
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Error("Handler error").
		Int("status", code).
		Err(err).
		Send()

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
		"error":   nil,
	})
}
