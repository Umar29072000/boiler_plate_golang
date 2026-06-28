package services

import (
	"boiler_plate_be_golang/internal/config"
	"boiler_plate_be_golang/internal/models"
	"boiler_plate_be_golang/internal/repositories"
	"boiler_plate_be_golang/pkg/email"
	"boiler_plate_be_golang/pkg/utils"
	"boiler_plate_be_golang/pkg/validator"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     *repositories.UserRepository
	emailService *email.EmailService
}

// NewAuthService creates new auth service
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		emailService: email.NewEmailService(),
	}
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

// Register registers new user and sends verification email
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

	// Generate email verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	// Set verification expiry (24 hours)
	verificationExpiry := time.Now().Add(24 * time.Hour)

	// Create user
	user := models.User{
		Name:                     req.Name,
		Email:                    req.Email,
		Password:                 hashedPassword,
		Role:                     "user",
		IsEmailVerified:          false,
		EmailVerificationToken:   &verificationToken,
		EmailVerificationExpires: &verificationExpiry,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Generate verification URL
	verificationURL := fmt.Sprintf("%s/api/auth/verify-email/%s", config.App.App.URL, verificationToken)

	// Send welcome email with verification link (async)
	go func() {
		if err := s.emailService.SendWelcomeEmail(user.Name, user.Email, verificationURL); err != nil {
			// Log error but don't fail registration
			fmt.Printf("Failed to send welcome email: %v\n", err)
		}
	}()

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

// Login authenticates user and updates last login
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

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(user); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
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

// VerifyEmail verifies user email with token
func (s *AuthService) VerifyEmail(token string) error {
	// Find user by verification token
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invalid or expired verification token")
		}
		return err
	}

	// Check if token is expired
	if user.EmailVerificationExpires != nil && time.Now().After(*user.EmailVerificationExpires) {
		return errors.New("verification token has expired")
	}

	// Check if already verified
	if user.IsEmailVerified {
		return errors.New("email already verified")
	}

	// Update user
	user.IsEmailVerified = true
	user.EmailVerificationToken = nil
	user.EmailVerificationExpires = nil

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// ResendVerificationEmail resends verification email
func (s *AuthService) ResendVerificationEmail(email string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	// Check if already verified
	if user.IsEmailVerified {
		return errors.New("email already verified")
	}

	// Generate new verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return err
	}

	// Set verification expiry (24 hours)
	verificationExpiry := time.Now().Add(24 * time.Hour)

	// Update user
	user.EmailVerificationToken = &verificationToken
	user.EmailVerificationExpires = &verificationExpiry

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Generate verification URL
	verificationURL := fmt.Sprintf("%s/api/auth/verify-email/%s", config.App.App.URL, verificationToken)

	// Send verification email
	if err := s.emailService.SendVerificationEmail(user.Name, user.Email, verificationURL); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// ForgotPassword sends password reset email
func (s *AuthService) ForgotPassword(email string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if user exists or not for security
			return nil
		}
		return err
	}

	// Generate password reset token
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	// Set reset expiry (15 minutes)
	resetExpiry := time.Now().Add(15 * time.Minute)

	// Update user
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpires = &resetExpiry

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Generate reset URL
	resetURL := fmt.Sprintf("%s/api/auth/reset-password/%s", config.App.App.URL, resetToken)

	// Send password reset email
	if err := s.emailService.SendPasswordResetEmail(user.Name, user.Email, resetURL); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate password
	if !validator.ValidatePassword(newPassword) {
		return errors.New("password must be at least 8 characters")
	}

	// Find user by reset token
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invalid or expired reset token")
		}
		return err
	}

	// Check if token is expired
	if user.PasswordResetExpires != nil && time.Now().After(*user.PasswordResetExpires) {
		return errors.New("reset token has expired")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user
	user.Password = hashedPassword
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Send password changed confirmation email (async)
	go func() {
		if err := s.emailService.SendPasswordChangedEmail(user.Name, user.Email); err != nil {
			fmt.Printf("Failed to send password changed email: %v\n", err)
		}
	}()

	return nil
}
