package domains

import "time"

// User represents user domain model
type User struct {
	ID                       string     `json:"id"`
	Name                     string     `json:"name"`
	Email                    string     `json:"email"`
	Password                 string     `json:"-"`
	Role                     string     `json:"role"`
	IsEmailVerified          bool       `json:"isEmailVerified"`
	EmailVerificationToken   *string    `json:"-"`
	EmailVerificationExpires *time.Time `json:"-"`
	PasswordResetToken       *string    `json:"-"`
	PasswordResetExpires     *time.Time `json:"-"`
	LastLogin                *time.Time `json:"lastLogin,omitempty"`
	CreatedAt                time.Time  `json:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
	DeletedAt                *time.Time `json:"deletedAt,omitempty"`
}

// UserResponse is the response structure for user data
type UserResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	IsEmailVerified bool       `json:"isEmailVerified"`
	LastLogin       *time.Time `json:"lastLogin,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
