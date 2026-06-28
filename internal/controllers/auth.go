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

	return utils.SuccessResponse(ctx, fiber.StatusCreated, "User registered successfully. Please check your email to verify your account.", result)
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

// VerifyEmail handles email verification
// @route GET /api/auth/verify-email/:token
func (c *AuthController) VerifyEmail(ctx *fiber.Ctx) error {
	token := ctx.Params("token")

	if token == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Verification token is required", nil)
	}

	if err := c.authService.VerifyEmail(token); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Email verified successfully", nil)
}

// ResendVerificationEmail handles resending verification email
// @route POST /api/auth/resend-verification
func (c *AuthController) ResendVerificationEmail(ctx *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.Email == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Email is required", nil)
	}

	if err := c.authService.ResendVerificationEmail(req.Email); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Verification email sent successfully", nil)
}

// ForgotPassword handles forgot password request
// @route POST /api/auth/forgot-password
func (c *AuthController) ForgotPassword(ctx *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.Email == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Email is required", nil)
	}

	if err := c.authService.ForgotPassword(req.Email); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send password reset email", nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "If your email exists in our system, you will receive a password reset link", nil)
}

// ResetPassword handles password reset with token
// @route POST /api/auth/reset-password/:token
func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	token := ctx.Params("token")

	if token == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Reset token is required", nil)
	}

	var req struct {
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.Password == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Password is required", nil)
	}

	if err := c.authService.ResetPassword(token, req.Password); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Password reset successfully", nil)
}
