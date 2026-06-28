package service

import (
	"context"
	"net/http"

	"boiler_plate_be_golang/domains"
	"boiler_plate_be_golang/internal/repository"
	"boiler_plate_be_golang/internal/rest/response"
	"boiler_plate_be_golang/pkg/logger"
	model "boiler_plate_be_golang/pkg/model/database"

	"gorm.io/gorm"
)

// IUserService defines user service interface
type IUserService interface {
	GetProfile(ctx context.Context, userID string) (response.BaseResponse[domains.UserResponse], error)
	GetAllUsers(ctx context.Context, page, limit int) (response.PaginationResp[domains.UserResponse], error)
	UpdateProfile(ctx context.Context, userID string, name string) (response.BaseResponse[domains.UserResponse], error)
	DeleteUser(ctx context.Context, userID string) (response.BaseResponse[interface{}], error)
}

// UserService handles user business logic
type UserService struct {
	UserRepository repository.IUserRepository
}

// NewUserService creates new user service
func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{
		UserRepository: userRepo,
	}
}

// GetProfile gets user profile by ID
func (s *UserService) GetProfile(ctx context.Context, userID string) (response.BaseResponse[domains.UserResponse], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "UserService.GetProfile")

	user, err := s.UserRepository.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			loggerCtx.Errorf("User not found: %s", userID)
			return response.BaseResponse[domains.UserResponse]{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
				Data:    domains.UserResponse{},
			}, nil
		}
		loggerCtx.Errorf("Error finding user: %v", err)
		return response.BaseResponse[domains.UserResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.UserResponse{},
		}, nil
	}

	userResponse := mapToUserResponse(user)

	return response.BaseResponse[domains.UserResponse]{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    userResponse,
	}, nil
}

// GetAllUsers gets all users with pagination
func (s *UserService) GetAllUsers(ctx context.Context, page, limit int) (response.PaginationResp[domains.UserResponse], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "UserService.GetAllUsers")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := s.UserRepository.FindAll(ctx, limit, offset)
	if err != nil {
		loggerCtx.Errorf("Error finding users: %v", err)
		return response.PaginationResp[domains.UserResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    []domains.UserResponse{},
			Meta: response.Meta{
				Page:  page,
				Limit: limit,
				Total: 0,
			},
		}, nil
	}

	total, err := s.UserRepository.Count(ctx)
	if err != nil {
		loggerCtx.Errorf("Error counting users: %v", err)
		return response.PaginationResp[domains.UserResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    []domains.UserResponse{},
			Meta: response.Meta{
				Page:  page,
				Limit: limit,
				Total: 0,
			},
		}, nil
	}

	var responses []domains.UserResponse
	for _, user := range users {
		responses = append(responses, mapToUserResponse(&user))
	}

	return response.PaginationResp[domains.UserResponse]{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    responses,
		Meta: response.Meta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID string, name string) (response.BaseResponse[domains.UserResponse], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "UserService.UpdateProfile")

	user, err := s.UserRepository.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			loggerCtx.Errorf("User not found: %s", userID)
			return response.BaseResponse[domains.UserResponse]{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
				Data:    domains.UserResponse{},
			}, nil
		}
		loggerCtx.Errorf("Error finding user: %v", err)
		return response.BaseResponse[domains.UserResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.UserResponse{},
		}, nil
	}

	user.Name = name
	if err := s.UserRepository.Update(ctx, user); err != nil {
		loggerCtx.Errorf("Error updating user: %v", err)
		return response.BaseResponse[domains.UserResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.UserResponse{},
		}, nil
	}

	userResponse := mapToUserResponse(user)

	return response.BaseResponse[domains.UserResponse]{
		Code:    http.StatusOK,
		Message: "PROFILE_UPDATED",
		Data:    userResponse,
	}, nil
}

// DeleteUser deletes user
func (s *UserService) DeleteUser(ctx context.Context, userID string) (response.BaseResponse[interface{}], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "UserService.DeleteUser")

	_, err := s.UserRepository.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			loggerCtx.Errorf("User not found: %s", userID)
			return response.BaseResponse[interface{}]{
				Code:    http.StatusNotFound,
				Message: "USER_NOT_FOUND",
				Data:    nil,
			}, nil
		}
		loggerCtx.Errorf("Error finding user: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		}, nil
	}

	if err := s.UserRepository.Delete(ctx, userID); err != nil {
		loggerCtx.Errorf("Error deleting user: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		}, nil
	}

	return response.BaseResponse[interface{}]{
		Code:    http.StatusOK,
		Message: "USER_DELETED",
		Data:    nil,
	}, nil
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
