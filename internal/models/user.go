package models

import "time"

// User represents user model
type User struct {
	BaseModel
	Name                     string     `gorm:"type:varchar(255);not null" json:"name"`
	Email                    string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password                 string     `gorm:"type:varchar(255);not null" json:"-"`
	Role                     string     `gorm:"type:varchar(50);default:'user'" json:"role"`
	IsEmailVerified          bool       `gorm:"default:false" json:"isEmailVerified"`
	EmailVerificationToken   *string    `gorm:"type:varchar(255)" json:"-"`
	EmailVerificationExpires *time.Time `json:"-"`
	PasswordResetToken       *string    `gorm:"type:varchar(255)" json:"-"`
	PasswordResetExpires     *time.Time `json:"-"`
	LastLogin                *time.Time `json:"lastLogin,omitempty"`
}

// TableName overrides default table name
func (User) TableName() string {
	return "users"
}

// UserResponse is the response structure for user data
type UserResponse struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	IsEmailVerified bool       `json:"isEmailVerified"`
	LastLogin       *time.Time `json:"lastLogin,omitempty"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Role:            u.Role,
		IsEmailVerified: u.IsEmailVerified,
		LastLogin:       u.LastLogin,
	}
}
