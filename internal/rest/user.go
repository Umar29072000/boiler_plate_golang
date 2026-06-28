package rest

import (
	"net/http"
	"strconv"
	"strings"

	"boiler_plate_be_golang/domains/dto"
	"boiler_plate_be_golang/internal/rest/request"
	"boiler_plate_be_golang/internal/rest/response"
	"boiler_plate_be_golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// userHandler handles user requests
type userHandler struct {
	UserService service.IUserService
}

// InitUserHandler initializes user routes
func InitUserHandler(e fiber.Router, userService service.IUserService) {
	handler := &userHandler{
		UserService: userService,
	}

	userGroup := e.Group("/users")
	userGroup.Get("/profile", handler.GetProfile)
	userGroup.Get("", handler.GetAllUsers)
	userGroup.Put("/profile", handler.UpdateProfile)
	userGroup.Delete("/:id", handler.DeleteUser)
}

// GetProfile handles get user profile
func (h *userHandler) GetProfile(c *fiber.Ctx) error {
	var tag string = "internal.rest.user.GetProfile."

	userID := c.Locals("user_id")
	if userID == nil {
		logrus.WithFields(logrus.Fields{
			"tag": tag + "01",
		}).Error("user_id not found in context")

		return c.Status(http.StatusUnauthorized).JSON(response.BaseResponse{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    nil,
		})
	}

	user, err := h.UserService.FindByID(c.Context(), userID.(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to get user profile (from user service)")

		switch err.Error() {
		case "USER_NOT_FOUND":
			return c.Status(http.StatusNotFound).JSON(response.BaseResponse{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
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
		Message: "SUCCESS",
		Data:    user,
	})
}

// GetAllUsers handles get all users with pagination
func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.user.GetAllUsers."
		req request.GetUserRequest
	)

	// Bind query params
	if err := c.QueryParser(&req); err != nil {
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

	// Convert page and limit from string to int
	var page, limit int
	var err error

	if req.Page != "" {
		page, err = strconv.Atoi(req.Page)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "03",
				"error": err.Error(),
			}).Error("failed to convert page from string to int")

			return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    http.StatusInternalServerError,
				Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				Data:    nil,
			})
		}
	}

	if page == 0 {
		page = 1
	}

	if req.Limit != "" {
		limit, err = strconv.Atoi(req.Limit)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "04",
				"error": err.Error(),
			}).Error("failed to convert limit from string to int")

			return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    http.StatusInternalServerError,
				Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				Data:    nil,
			})
		}
	}

	if limit == 0 {
		limit = 10
	}

	// Call service
	userData, err := h.UserService.Show(c.Context(), dto.GetUserData{
		Page:                  page,
		Limit:                 limit,
		Field:                 req.Field,
		Sort:                  req.Sort,
		Search:                req.Search,
		DisableCalculateTotal: req.DisableCalculateTotal,
		ID:                    req.ID,
		Email:                 req.Email,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "05",
			"error": err.Error(),
		}).Error("failed to get users (from user service)")

		return c.Status(http.StatusInternalServerError).JSON(response.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    userData,
	})
}

// UpdateProfile handles update user profile
func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.user.UpdateProfile."
		req request.UpdateProfileRequest
	)

	userID := c.Locals("user_id")
	if userID == nil {
		logrus.WithFields(logrus.Fields{
			"tag": tag + "01",
		}).Error("user_id not found in context")

		return c.Status(http.StatusUnauthorized).JSON(response.BaseResponse{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    nil,
		})
	}

	// Bind request body
	if err := c.BodyParser(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
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
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("invalid validation")

		return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VALIDATION",
			Data:    err,
		})
	}

	// Get current user to preserve other fields
	currentUser, err := h.UserService.FindByID(c.Context(), userID.(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": err.Error(),
		}).Error("failed to find user")

		switch err.Error() {
		case "USER_NOT_FOUND":
			return c.Status(http.StatusNotFound).JSON(response.BaseResponse{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
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

	// Update user
	userData, err := h.UserService.Update(c.Context(), dto.UpdateUserData{
		ID:    userID.(string),
		Name:  req.Name,
		Email: currentUser.Email,
		Role:  currentUser.Role,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "05",
			"error": err.Error(),
		}).Error("failed to update user profile (from user service)")

		switch err.Error() {
		case "USER_NOT_FOUND":
			return c.Status(http.StatusNotFound).JSON(response.BaseResponse{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
				Data:    nil,
			})
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

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    userData,
	})
}

// DeleteUser handles delete user
func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	var (
		tag string = "internal.rest.user.DeleteUser."
		req request.DeleteUserRequest
	)

	// Bind path params
	req.ID = c.Params("id")

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

	if err := h.UserService.Delete(c.Context(), req.ID); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to delete user (from user service)")

		switch err.Error() {
		case "USER_NOT_FOUND":
			return c.Status(http.StatusNotFound).JSON(response.BaseResponse{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
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
		Message: "SUCCESS",
		Data:    nil,
	})
}
