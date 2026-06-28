package rest

import (
	"net/http"
	"strings"

	"boiler_plate_be_golang/domains"
	"boiler_plate_be_golang/internal/rest/request"
	"boiler_plate_be_golang/internal/rest/response"
	"boiler_plate_be_golang/internal/service"
	model "boiler_plate_be_golang/pkg/model/database"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// authHandler handles authentication requests
type authHandler struct {
	AuthService service.IAuthService
}

// InitAuthHandler initializes auth routes
func InitAuthHandler(e fiber.Router, authService service.IAuthService) {
	handler := &authHandler{
		AuthService: authService,
	}

	authGroup := e.Group("/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Get("/verify/:token", handler.VerifyEmail)
	authGroup.Post("/forgot-password", handler.ForgotPassword)
	authGroup.Post("/reset-password", handler.ResetPassword)
}

// Register handles user registration
func (h *authHandler) Register(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.auth.Register."
		req request.RegisterRequest
	)

	// Bind request body
	if err := c.BodyParser(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("bad request")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:    nil,
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("invalid validation")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VALIDATION",
			Data:    err,
		})
	}

	// Call service
	user, token, err := h.AuthService.Register(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to register user (from auth service)")

		switch err.Error() {
		case "EMAIL_FOUND":
			return c.Status(http.StatusConflict).JSON(response.BaseResponse{
				Code:    http.StatusConflict,
				Message: "EMAIL_ALREADY_EXISTS",
				Data:    nil,
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    http.StatusInternalServerError,
				Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				Data:    nil,
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(response.BaseResponse{
		Code:    http.StatusCreated,
		Message: "USER_REGISTERED",
		Data: map[string]interface{}{
			"user":  mapToUserResponse(user),
			"token": token,
		},
	})
}

// Login handles user login
func (h *authHandler) Login(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.auth.Login."
		req request.LoginRequest
	)

	// Bind request body
	if err := c.BodyParser(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("bad request")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:    nil,
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("invalid validation")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VALIDATION",
			Data:    err,
		})
	}

	// Call service
	user, token, err := h.AuthService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to login (from auth service)")

		switch err.Error() {
		case "INVALID_CREDENTIALS":
			return c.Status(http.StatusUnauthorized).JSON(response.BaseResponse{
				Code:    http.StatusUnauthorized,
				Message: "INVALID_CREDENTIALS",
				Data:    nil,
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    http.StatusInternalServerError,
				Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				Data:    nil,
			})
		}
	}

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "LOGIN_SUCCESS",
		Data: map[string]interface{}{
			"user":  mapToUserResponse(user),
			"token": token,
		},
	})
}



// VerifyEmail handles email verification
func (h *authHandler) VerifyEmail(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.auth.VerifyEmail."
		req request.VerifyEmailRequest
	)

	// Bind path params
	req.Token = c.Params("token")

	// Validate request
	if err := req.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("invalid validation")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VALIDATION",
			Data:    err,
		})
	}

	// Call service
	if err := h.AuthService.VerifyEmail(c.Context(), req.Token); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to verify email (from auth service)")

		switch err.Error() {
		case "INVALID_VERIFICATION_TOKEN":
			return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
				Code:    http.StatusBadRequest,
				Message: "INVALID_VERIFICATION_TOKEN",
				Data:    nil,
			})
		case "VERIFICATION_TOKEN_EXPIRED":
			return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
				Code:    http.StatusBadRequest,
				Message: "VERIFICATION_TOKEN_EXPIRED",
				Data:    nil,
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    http.StatusInternalServerError,
				Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				Data:    nil,
			})
		}
	}

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "EMAIL_VERIFIED",
		Data:    nil,
	})
}

// ForgotPassword handles forgot password request
func (h *authHandler) ForgotPassword(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.auth.ForgotPassword."
		req request.ForgotPasswordRequest
	)

	// Bind request body
	if err := c.BodyParser(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("bad request")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:    nil,
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("invalid validation")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VALIDATION",
			Data:    err,
		})
	}

	// Call service
	if err := h.AuthService.ForgotPassword(c.Context(), req.Email); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to process forgot password (from auth service)")

		return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "PASSWORD_RESET_EMAIL_SENT",
		Data:    nil,
	})
}

// ResetPassword handles password reset
func (h *authHandler) ResetPassword(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.auth.ResetPassword."
		req request.ResetPasswordRequest
	)

	// Bind request body
	if err := c.BodyParser(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("bad request")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:    nil,
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("invalid validation")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VALIDATION",
			Data:    err,
		})
	}

	// Call service
	if err := h.AuthService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to reset password (from auth service)")

		switch err.Error() {
		case "INVALID_RESET_TOKEN":
			return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
				Code:    http.StatusBadRequest,
				Message: "INVALID_RESET_TOKEN",
				Data:    nil,
			})
		case "RESET_TOKEN_EXPIRED":
			return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
				Code:    http.StatusBadRequest,
				Message: "RESET_TOKEN_EXPIRED",
				Data:    nil,
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    http.StatusInternalServerError,
				Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				Data:    nil,
			})
		}
	}

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "PASSWORD_RESET_SUCCESS",
		Data:    nil,
	})
}

// mapToUserResponse maps database model to domain response
func mapToUserResponse(user *model.User) domains.UserResponse {
	return domains.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		LastLogin:       user.LastLogin,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
