package controllers

import (
	"boiler_plate_be_golang/internal/services"
	"boiler_plate_be_golang/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserController handles user requests
type UserController struct {
	userService *services.UserService
}

// NewUserController creates new user controller
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile handles get user profile
// @route GET /api/users/profile
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	result, err := c.userService.GetProfile(userID)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Profile retrieved successfully", result)
}

// GetAllUsers handles get all users with pagination
// @route GET /api/users
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := c.userService.GetAllUsers(page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve users", err.Error())
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Users retrieved successfully", fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateProfile handles update user profile
// @route PUT /api/users/profile
func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	var req struct {
		Name string `json:"name"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	result, err := c.userService.UpdateProfile(userID, req.Name)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Profile updated successfully", result)
}

// DeleteUser handles delete user
// @route DELETE /api/users/:id
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid user ID", err.Error())
	}

	if err := c.userService.DeleteUser(uint(id)); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "User deleted successfully", nil)
}
