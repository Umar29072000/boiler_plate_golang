package services

import (
	"boiler_plate_be_golang/internal/models"
	"boiler_plate_be_golang/internal/repositories"
	"boiler_plate_be_golang/pkg/utils"
	"boiler_plate_be_golang/pkg/validator"
	"errors"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService creates new auth service
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

// Register registers new user
func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Validate required fields
	validationErrors := validator.ValidateRequired(map[string]string{
		"name":     req.Name,
		"email":    req.Email,
		"password": req.Password,
	})
	if len(validationErrors) > 0 {
		return nil, errors.New(validationErrors[0])
	}

	// Validate email
	if !validator.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password
	if !validator.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if user exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates user
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Validate required fields
	validationErrors := validator.ValidateRequired(map[string]string{
		"email":    req.Email,
		"password": req.Password,
	})
	if len(validationErrors) > 0 {
		return nil, errors.New(validationErrors[0])
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Compare password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}
