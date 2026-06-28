package service

import (
	"context"
	"errors"
	"time"

	"boiler_plate_be_golang/app/config"
	"boiler_plate_be_golang/domains/dto"
	"boiler_plate_be_golang/internal/repository"
	model "boiler_plate_be_golang/pkg/model/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// IAuthService defines auth service interface
type IAuthService interface {
	Register(ctx context.Context, name, email, password string) (user *model.User, token string, err error)
	Login(ctx context.Context, email, password string) (user *model.User, token string, err error)
	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

// AuthService handles authentication business logic
type AuthService struct {
	UserRepository repository.IUserRepository
	JWTConfig      config.JWT
}

// NewAuthService creates new auth service
func NewAuthService(userRepo repository.IUserRepository, jwtConfig config.JWT) *AuthService {
	return &AuthService{
		UserRepository: userRepo,
		JWTConfig:      jwtConfig,
	}
}

// Register handles user registration
func (s *AuthService) Register(ctx context.Context, name, email, password string) (user *model.User, token string, err error) {
	var tag string = "internal.service.auth.Register."

	// Check if user already exists
	existingUser, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to check existing user (from user repository)")
		return nil, "", err
	}

	if existingUser != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"email": email,
		}).Error("email already exists")
		return nil, "", errors.New("EMAIL_FOUND")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to hash password")
		return nil, "", err
	}

	// Generate verification token
	verificationToken := uuid.New().String()
	verificationExpires := time.Now().Add(24 * time.Hour)

	// Create user using repository DTO
	user, err = s.UserRepository.Create(ctx, dto.CreateUserData{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": err.Error(),
		}).Error("failed to create user (from user repository)")
		return nil, "", err
	}

	// Update user with verification token
	user.EmailVerificationToken = &verificationToken
	user.EmailVerificationExpires = &verificationExpires
	user.IsEmailVerified = false

	updateReq := dto.UpdateUserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	
	if _, err = s.UserRepository.Update(ctx, updateReq); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "05",
			"error": err.Error(),
		}).Error("failed to update user with verification token")
	}

	// Generate token (7 days expiration)
	token, err = s.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "06",
			"error": err.Error(),
		}).Error("failed to generate token")
		return nil, "", err
	}

	// TODO: Send verification email
	logrus.WithFields(logrus.Fields{
		"tag":   tag + "07",
		"email": user.Email,
		"token": verificationToken,
	}).Info("user registered successfully")

	return user, token, nil
}

// Login handles user login
func (s *AuthService) Login(ctx context.Context, email, password string) (user *model.User, token string, err error) {
	var tag string = "internal.service.auth.Login."

	// Find user by email
	user, err = s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "01",
				"email": email,
			}).Error("user not found")
			return nil, "", errors.New("INVALID_CREDENTIALS")
		}
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to find user (from user repository)")
		return nil, "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"email": email,
		}).Error("invalid password")
		return nil, "", errors.New("INVALID_CREDENTIALS")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	updateReq := dto.UpdateUserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	if _, err := s.UserRepository.Update(ctx, updateReq); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": err.Error(),
		}).Error("failed to update last login")
	}

	// Generate token (7 days expiration)
	token, err = s.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "05",
			"error": err.Error(),
		}).Error("failed to generate token")
		return nil, "", err
	}

	return user, token, nil
}



// VerifyEmail handles email verification
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	var tag string = "internal.service.auth.VerifyEmail."

	user, err := s.UserRepository.FindByVerificationToken(ctx, token)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("invalid verification token")
		return errors.New("INVALID_VERIFICATION_TOKEN")
	}

	// Check if token is expired
	if user.EmailVerificationExpires != nil && user.EmailVerificationExpires.Before(time.Now()) {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"email": user.Email,
		}).Error("verification token expired")
		return errors.New("VERIFICATION_TOKEN_EXPIRED")
	}

	// Update user
	user.IsEmailVerified = true
	user.EmailVerificationToken = nil
	user.EmailVerificationExpires = nil

	updateReq := dto.UpdateUserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	if _, err := s.UserRepository.Update(ctx, updateReq); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to update user")
		return err
	}

	return nil
}

// ForgotPassword handles forgot password request
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	var tag string = "internal.service.auth.ForgotPassword."

	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists or not - return success anyway
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "01",
				"email": email,
			}).Info("password reset requested for non-existent email")
			return nil
		}
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to find user")
		return err
	}

	// Generate reset token
	resetToken := uuid.New().String()
	resetExpires := time.Now().Add(1 * time.Hour)

	user.PasswordResetToken = &resetToken
	user.PasswordResetExpires = &resetExpires

	updateReq := dto.UpdateUserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	if _, err := s.UserRepository.Update(ctx, updateReq); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to update user with reset token")
		return err
	}

	// TODO: Send password reset email
	logrus.WithFields(logrus.Fields{
		"tag":   tag + "04",
		"email": user.Email,
		"token": resetToken,
	}).Info("password reset email should be sent")

	return nil
}

// ResetPassword handles password reset
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	var tag string = "internal.service.auth.ResetPassword."

	user, err := s.UserRepository.FindByResetToken(ctx, token)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("invalid reset token")
		return errors.New("INVALID_RESET_TOKEN")
	}

	// Check if token is expired
	if user.PasswordResetExpires != nil && user.PasswordResetExpires.Before(time.Now()) {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"email": user.Email,
		}).Error("reset token expired")
		return errors.New("RESET_TOKEN_EXPIRED")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to hash password")
		return err
	}

	// Update user
	user.Password = string(hashedPassword)
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil

	updateReq := dto.UpdateUserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	if _, err := s.UserRepository.Update(ctx, updateReq); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": err.Error(),
		}).Error("failed to update user password")
		return err
	}

	return nil
}

// generateToken generates JWT token with 7 days expiration
func (s *AuthService) generateToken(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTConfig.Secret))
}
