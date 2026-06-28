package service

import (
	"context"
	"errors"

	"boiler_plate_be_golang/domains/dto"
	"boiler_plate_be_golang/internal/repository"
	"boiler_plate_be_golang/internal/rest/response"
	model "boiler_plate_be_golang/pkg/model/database"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IUserService defines user service interface
type IUserService interface {
	Show(ctx context.Context, req dto.GetUserData) (res response.PaginationResponse, err error)
	FindByID(ctx context.Context, id string) (user *model.User, err error)
	Update(ctx context.Context, req dto.UpdateUserData) (user *model.User, err error)
	Delete(ctx context.Context, id string) error
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

// Show gets users with pagination
func (s *UserService) Show(ctx context.Context, req dto.GetUserData) (res response.PaginationResponse, err error) {
	var tag string = "internal.service.user.Show."

	data, err := s.UserRepository.Show(ctx, req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to get users (from user repository)")
		return data, err
	}

	return data, nil
}

// FindByID finds user by ID
func (s *UserService) FindByID(ctx context.Context, id string) (user *model.User, err error) {
	var tag string = "internal.service.user.FindByID."

	user, err = s.UserRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{
				"tag":    tag + "01",
				"userID": id,
			}).Error("user not found")
			return nil, errors.New("USER_NOT_FOUND")
		}

		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to find user by ID (from user repository)")
		return nil, err
	}

	return user, nil
}

// Update updates user profile
func (s *UserService) Update(ctx context.Context, req dto.UpdateUserData) (user *model.User, err error) {
	var tag string = "internal.service.user.Update."

	// Check if user exists
	existingUser, err := s.UserRepository.FindByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{
				"tag":    tag + "01",
				"userID": req.ID,
			}).Error("user not found")
			return nil, errors.New("USER_NOT_FOUND")
		}

		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to fetch current user data (from user repository)")
		return nil, err
	}

	// Check if email is being changed and already exists
	if req.Email != "" && req.Email != existingUser.Email {
		emailUser, err := s.UserRepository.FindByEmail(ctx, req.Email)
		if err == nil && emailUser.ID != req.ID {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "03",
				"email": req.Email,
			}).Error("email already exists")
			return nil, errors.New("EMAIL_FOUND")
		}
	}

	user, err = s.UserRepository.Update(ctx, req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": err.Error(),
		}).Error("failed to update user (from user repository)")
		return nil, err
	}

	return user, nil
}

// Delete deletes user
func (s *UserService) Delete(ctx context.Context, id string) error {
	var tag string = "internal.service.user.Delete."

	// Check if user exists first
	_, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{
				"tag":    tag + "01",
				"userID": id,
			}).Error("user not found")
			return errors.New("USER_NOT_FOUND")
		}

		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to find user (from user repository)")
		return err
	}

	if err := s.UserRepository.Delete(ctx, id); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to delete user (from user repository)")
		return err
	}

	return nil
}
