package rest

import (
	"net/http"
	"strconv"

	"boiler_plate_be_golang/internal/service"
	"boiler_plate_be_golang/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user requests
type UserHandler struct {
	UserService service.IUserService
}

// InitUserHandler initializes user routes
func InitUserHandler(e fiber.Router, userService service.IUserService) {
	handler := &UserHandler{
		UserService: userService,
	}

	userGroup := e.Group("/users")
	userGroup.Get("/profile", handler.GetProfile)
	userGroup.Get("", handler.GetAllUsers)
	userGroup.Put("/profile", handler.UpdateProfile)
	userGroup.Delete("/:id", handler.DeleteUser)
}

// GetProfile handles get user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "UserHandler.GetProfile")

	userID := c.Locals("user_id")
	if userID == nil {
		loggerCtx.Error("user_id not found in context")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "UNAUTHORIZED",
			"data":    nil,
		})
	}

	res, err := h.UserService.GetProfile(ctx, userID.(string))
	if err != nil {
		loggerCtx.Errorf("Error getting profile: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// GetAllUsers handles get all users with pagination
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "UserHandler.GetAllUsers")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, err := h.UserService.GetAllUsers(ctx, page, limit)
	if err != nil {
		loggerCtx.Errorf("Error getting users: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// UpdateProfile handles update user profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "UserHandler.UpdateProfile")

	userID := c.Locals("user_id")
	if userID == nil {
		loggerCtx.Error("user_id not found in context")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "UNAUTHORIZED",
			"data":    nil,
		})
	}

	var req struct {
		Name string `json:"name" validate:"required,min=3,max=100"`
	}

	if err := c.BodyParser(&req); err != nil {
		loggerCtx.Errorf("Error parsing request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.UserService.UpdateProfile(ctx, userID.(string), req.Name)
	if err != nil {
		loggerCtx.Errorf("Error updating profile: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}

// DeleteUser handles delete user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()
	loggerCtx := logger.GetLoggerContext(ctx, "UserHandler.DeleteUser")

	id := c.Params("id")
	if id == "" {
		loggerCtx.Error("id parameter is missing")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "INVALID_REQUEST",
			"data":    nil,
		})
	}

	res, err := h.UserService.DeleteUser(ctx, id)
	if err != nil {
		loggerCtx.Errorf("Error deleting user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "INTERNAL_SERVER_ERROR",
			"data":    nil,
		})
	}

	return c.Status(res.Code).JSON(res)
}
