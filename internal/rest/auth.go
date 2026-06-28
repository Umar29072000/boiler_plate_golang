package rest

import (
	"net/http"

	"boiler_plate_be_golang/domains"
	"boiler_plate_be_golang/internal/service"
	"boiler_plate_be_golang/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	AuthService service.IAuthService
}

// InitAuthHandler initializes auth routes
func InitAuthHandler(e fiber.Router, authService service.IAuthService) {
	handler := &AuthHandler{
		AuthService: authService,
	}

	authGroup := e.Group("/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.RefreshToken)
	authGroup.Get("/verify/:token", handler.VerifyEmail)
	authGroup.Post("/forgot-password", handler.ForgotPassword)
	authGroup.Post("/reset-password", handler.ResetPassword)
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "AuthHandler.Register")

	var req domains.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		loggerCtx.Errorf("Error parsing request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	// TODO: Add validation
	res, err := h.AuthService.Register(ctx, req)
	if err != nil {
		loggerCtx.Errorf("Error registering user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "AuthHandler.Login")

	var req domains.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		loggerCtx.Errorf("Error parsing request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.AuthService.Login(ctx, req)
	if err != nil {
		loggerCtx.Errorf("Error logging in: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "AuthHandler.RefreshToken")

	var req domains.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		loggerCtx.Errorf("Error parsing request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.AuthService.RefreshToken(ctx, req)
	if err != nil {
		loggerCtx.Errorf("Error refreshing token: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "AuthHandler.VerifyEmail")

	token := c.Params("token")
	if token == "" {
		loggerCtx.Error("token parameter is missing")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.AuthService.VerifyEmail(ctx, token)
	if err != nil {
		loggerCtx.Errorf("Error verifying email: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// ForgotPassword handles forgot password request
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "AuthHandler.ForgotPassword")

	var req domains.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		loggerCtx.Errorf("Error parsing request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.AuthService.ForgotPassword(ctx, req.Email)
	if err != nil {
		loggerCtx.Errorf("Error processing forgot password: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "AuthHandler.ResetPassword")

	var req domains.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		loggerCtx.Errorf("Error parsing request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.AuthService.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		loggerCtx.Errorf("Error resetting password: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}
