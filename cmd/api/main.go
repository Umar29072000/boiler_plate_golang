package main

import (
	"boiler_plate_be_golang/internal/config"
	"boiler_plate_be_golang/internal/database"
	"boiler_plate_be_golang/internal/database/migrations"
	"boiler_plate_be_golang/internal/middleware"
	"boiler_plate_be_golang/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run migrations
	if err := migrations.Migrate(database.GetDB()); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      config.App.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.ErrorHandler())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := config.App.App.Port
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", config.App.App.Env)
	log.Printf("URL: %s", config.App.App.URL)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
		"error":   nil,
	})
}
