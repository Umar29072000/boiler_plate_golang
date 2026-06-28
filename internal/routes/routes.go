package routes

import (
	"boiler_plate_be_golang/internal/controllers"
	"boiler_plate_be_golang/internal/database"
	"boiler_plate_be_golang/internal/middleware"
	"boiler_plate_be_golang/internal/repositories"
	"boiler_plate_be_golang/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// Initialize dependencies
	db := database.GetDB()
	
	// Repositories
	userRepo := repositories.NewUserRepository(db)
	
	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	
	// Controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	// API routes
	api := app.Group("/api")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Server is running",
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/verify-email/:token", authController.VerifyEmail)
	auth.Post("/resend-verification", authController.ResendVerificationEmail)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password/:token", authController.ResetPassword)

	// User routes (protected)
	users := api.Group("/users")
	users.Use(middleware.Auth())
	
	users.Get("/profile", userController.GetProfile)
	users.Put("/profile", userController.UpdateProfile)
	users.Get("/", middleware.AdminOnly(), userController.GetAllUsers)
	users.Delete("/:id", middleware.AdminOnly(), userController.DeleteUser)
}
