package domains

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}
