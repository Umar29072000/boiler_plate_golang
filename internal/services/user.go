package services

import (
	"boiler_plate_be_golang/internal/models"
	"boiler_plate_be_golang/internal/repositories"
	"errors"

	"gorm.io/gorm"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates new user service
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetProfile gets user profile by ID
func (s *UserService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetAllUsers gets all users with pagination
func (s *UserService) GetAllUsers(page, limit int) ([]models.UserResponse, int64, error) {
	offset := (page - 1) * limit

	users, err := s.userRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, total, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uint, name string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Name = name
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser deletes user
func (s *UserService) DeleteUser(userID uint) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(userID)
}
