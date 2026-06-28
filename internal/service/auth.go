package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"boiler_plate_be_golang/app/config"
	"boiler_plate_be_golang/domains"
	"boiler_plate_be_golang/internal/repository"
	"boiler_plate_be_golang/internal/rest/response"
	"boiler_plate_be_golang/pkg/logger"
	model "boiler_plate_be_golang/pkg/model/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// IAuthService defines auth service interface
type IAuthService interface {
	Register(ctx context.Context, req domains.RegisterRequest) (response.BaseResponse[domains.LoginResponse], error)
	Login(ctx context.Context, req domains.LoginRequest) (response.BaseResponse[domains.LoginResponse], error)
	RefreshToken(ctx context.Context, req domains.RefreshTokenRequest) (response.BaseResponse[domains.LoginResponse], error)
	VerifyEmail(ctx context.Context, token string) (response.BaseResponse[interface{}], error)
	ForgotPassword(ctx context.Context, email string) (response.BaseResponse[interface{}], error)
	ResetPassword(ctx context.Context, token, newPassword string) (response.BaseResponse[interface{}], error)
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
func (s *AuthService) Register(ctx context.Context, req domains.RegisterRequest) (response.BaseResponse[domains.LoginResponse], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "AuthService.Register")

	// Check if user already exists
	existingUser, err := s.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		loggerCtx.Errorf("Error checking existing user: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	if existingUser != nil {
		loggerCtx.Warnf("User already exists: %s", req.Email)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusConflict,
			Message: "EMAIL_ALREADY_EXISTS",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		loggerCtx.Errorf("Error hashing password: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Generate verification token
	verificationToken := uuid.New().String()
	verificationExpires := time.Now().Add(24 * time.Hour)

	// Create user
	user := &model.User{
		ID:                       uuid.New().String(),
		Name:                     req.Name,
		Email:                    req.Email,
		Password:                 string(hashedPassword),
		Role:                     "user",
		IsEmailVerified:          false,
		EmailVerificationToken:   &verificationToken,
		EmailVerificationExpires: &verificationExpires,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	if err := s.UserRepository.Create(ctx, user); err != nil {
		loggerCtx.Errorf("Error creating user: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		loggerCtx.Errorf("Error generating access token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		loggerCtx.Errorf("Error generating refresh token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// TODO: Send verification email
	loggerCtx.Infof("User registered successfully: %s (verification token: %s)", user.Email, verificationToken)

	return response.BaseResponse[domains.LoginResponse]{
		Code:    http.StatusCreated,
		Message: "USER_REGISTERED",
		Data: domains.LoginResponse{
			User:         mapToUserResponse(user),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// Login handles user login
func (s *AuthService) Login(ctx context.Context, req domains.LoginRequest) (response.BaseResponse[domains.LoginResponse], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "AuthService.Login")

	// Find user by email
	user, err := s.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			loggerCtx.Warnf("User not found: %s", req.Email)
			return response.BaseResponse[domains.LoginResponse]{
				Code:    http.StatusUnauthorized,
				Message: "INVALID_CREDENTIALS",
				Data:    domains.LoginResponse{},
			}, nil
		}
		loggerCtx.Errorf("Error finding user: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		loggerCtx.Warnf("Invalid password for user: %s", req.Email)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusUnauthorized,
			Message: "INVALID_CREDENTIALS",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.UserRepository.Update(ctx, user); err != nil {
		loggerCtx.Errorf("Error updating last login: %v", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		loggerCtx.Errorf("Error generating access token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		loggerCtx.Errorf("Error generating refresh token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	return response.BaseResponse[domains.LoginResponse]{
		Code:    http.StatusOK,
		Message: "LOGIN_SUCCESS",
		Data: domains.LoginResponse{
			User:         mapToUserResponse(user),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// RefreshToken handles token refresh
func (s *AuthService) RefreshToken(ctx context.Context, req domains.RefreshTokenRequest) (response.BaseResponse[domains.LoginResponse], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "AuthService.RefreshToken")

	// Parse and validate refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.JWTConfig.Secret), nil
	})

	if err != nil || !token.Valid {
		loggerCtx.Errorf("Invalid refresh token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusUnauthorized,
			Message: "INVALID_REFRESH_TOKEN",
			Data:    domains.LoginResponse{},
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		loggerCtx.Error("Invalid token claims")
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusUnauthorized,
			Message: "INVALID_REFRESH_TOKEN",
			Data:    domains.LoginResponse{},
		}, nil
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		loggerCtx.Error("Invalid user_id in token")
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusUnauthorized,
			Message: "INVALID_REFRESH_TOKEN",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Find user
	user, err := s.UserRepository.FindByID(ctx, userID)
	if err != nil {
		loggerCtx.Errorf("User not found: %s", userID)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusUnauthorized,
			Message: "USER_NOT_FOUND",
			Data:    domains.LoginResponse{},
		}, nil
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		loggerCtx.Errorf("Error generating access token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		loggerCtx.Errorf("Error generating refresh token: %v", err)
		return response.BaseResponse[domains.LoginResponse]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    domains.LoginResponse{},
		}, nil
	}

	return response.BaseResponse[domains.LoginResponse]{
		Code:    http.StatusOK,
		Message: "TOKEN_REFRESHED",
		Data: domains.LoginResponse{
			User:         mapToUserResponse(user),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// VerifyEmail handles email verification
func (s *AuthService) VerifyEmail(ctx context.Context, token string) (response.BaseResponse[interface{}], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "AuthService.VerifyEmail")

	user, err := s.UserRepository.FindByVerificationToken(ctx, token)
	if err != nil {
		loggerCtx.Errorf("Invalid verification token: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusBadRequest,
			Message: "INVALID_VERIFICATION_TOKEN",
			Data:    nil,
		}, nil
	}

	// Check if token is expired
	if user.EmailVerificationExpires != nil && user.EmailVerificationExpires.Before(time.Now()) {
		loggerCtx.Warnf("Verification token expired for user: %s", user.Email)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusBadRequest,
			Message: "VERIFICATION_TOKEN_EXPIRED",
			Data:    nil,
		}, nil
	}

	// Update user
	user.IsEmailVerified = true
	user.EmailVerificationToken = nil
	user.EmailVerificationExpires = nil

	if err := s.UserRepository.Update(ctx, user); err != nil {
		loggerCtx.Errorf("Error updating user: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		}, nil
	}

	return response.BaseResponse[interface{}]{
		Code:    http.StatusOK,
		Message: "EMAIL_VERIFIED",
		Data:    nil,
	}, nil
}

// ForgotPassword handles forgot password request
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (response.BaseResponse[interface{}], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "AuthService.ForgotPassword")

	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if email exists or not
			return response.BaseResponse[interface{}]{
				Code:    http.StatusOK,
				Message: "PASSWORD_RESET_EMAIL_SENT",
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

	// Generate reset token
	resetToken := uuid.New().String()
	resetExpires := time.Now().Add(1 * time.Hour)

	user.PasswordResetToken = &resetToken
	user.PasswordResetExpires = &resetExpires

	if err := s.UserRepository.Update(ctx, user); err != nil {
		loggerCtx.Errorf("Error updating user: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		}, nil
	}

	// TODO: Send password reset email
	loggerCtx.Infof("Password reset requested for: %s (reset token: %s)", user.Email, resetToken)

	return response.BaseResponse[interface{}]{
		Code:    http.StatusOK,
		Message: "PASSWORD_RESET_EMAIL_SENT",
		Data:    nil,
	}, nil
}

// ResetPassword handles password reset
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) (response.BaseResponse[interface{}], error) {
	loggerCtx := logger.GetLoggerContext(ctx, "AuthService.ResetPassword")

	user, err := s.UserRepository.FindByResetToken(ctx, token)
	if err != nil {
		loggerCtx.Errorf("Invalid reset token: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusBadRequest,
			Message: "INVALID_RESET_TOKEN",
			Data:    nil,
		}, nil
	}

	// Check if token is expired
	if user.PasswordResetExpires != nil && user.PasswordResetExpires.Before(time.Now()) {
		loggerCtx.Warnf("Reset token expired for user: %s", user.Email)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusBadRequest,
			Message: "RESET_TOKEN_EXPIRED",
			Data:    nil,
		}, nil
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		loggerCtx.Errorf("Error hashing password: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		}, nil
	}

	// Update user
	user.Password = string(hashedPassword)
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil

	if err := s.UserRepository.Update(ctx, user); err != nil {
		loggerCtx.Errorf("Error updating user: %v", err)
		return response.BaseResponse[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "INTERNAL_SERVER_ERROR",
			Data:    nil,
		}, nil
	}

	return response.BaseResponse[interface{}]{
		Code:    http.StatusOK,
		Message: "PASSWORD_RESET_SUCCESS",
		Data:    nil,
	}, nil
}

// generateAccessToken generates JWT access token
func (s *AuthService) generateAccessToken(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(s.JWTConfig.Expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTConfig.Secret))
}

// generateRefreshToken generates JWT refresh token
func (s *AuthService) generateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTConfig.Secret))
}
