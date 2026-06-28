package controllers

import (
	"boiler_plate_be_golang/internal/services"
	"boiler_plate_be_golang/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication requests
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates new auth controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register handles user registration
// @route POST /api/auth/register
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req services.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	result, err := c.authService.Register(req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusCreated, "User registered successfully", result)
}

// Login handles user login
// @route POST /api/auth/login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req services.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	result, err := c.authService.Login(req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Login successful", result)
}
