package middleware

import (
	"boiler_plate_be_golang/internal/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns configured CORS middleware
func CORS() fiber.Handler {
	allowedOrigins := config.App.CORS.AllowedOrigins
	
	return cors.New(cors.Config{
		AllowOrigins:     strings.ReplaceAll(allowedOrigins, " ", ""),
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	})
}
