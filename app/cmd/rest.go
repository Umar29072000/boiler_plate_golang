package cmd

import (
	"os"

	"boiler_plate_be_golang/middleware"
	restHandler "boiler_plate_be_golang/internal/rest"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	Run:   restServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}

func restServer(cmd *cobra.Command, args []string) {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      rootConfig.App.ServiceName,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.SecurityHeaders())
	app.Use(cors.New(cors.Config{
		AllowOrigins: rootConfig.CORS.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		Max:      rootConfig.RateLimit.Max,
		Duration: rootConfig.RateLimit.Duration,
	}))

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Initialize handlers
	restHandler.InitHealthHandler(v1)
	restHandler.InitAuthHandler(v1, authService)
	restHandler.InitUserHandler(v1, userService)

	// Start server
	port := rootConfig.App.Port
	logrus.Info("Server starting on port: ", port)
	logrus.Info("Environment: ", rootConfig.App.Env)
	logrus.Info("Service: ", rootConfig.App.ServiceName)

	if err := app.Listen(":" + port); err != nil {
		logrus.Error("Could not start server: ", err)
		os.Exit(1)
	}
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logrus.WithFields(logrus.Fields{
		"status": code,
		"error":  err.Error(),
		"path":   c.Path(),
		"method": c.Method(),
	}).Error("Handler error")

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": err.Error(),
		"data":    nil,
	})
}
