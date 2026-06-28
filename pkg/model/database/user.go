package database

import "time"

// User represents database user model for GORM
type User struct {
	ID                       string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name                     string     `gorm:"type:varchar(255);not null" json:"name"`
	Email                    string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password                 string     `gorm:"type:varchar(255);not null" json:"-"`
	Role                     string     `gorm:"type:varchar(50);default:'user'" json:"role"`
	IsEmailVerified          bool       `gorm:"default:false" json:"isEmailVerified"`
	EmailVerificationToken   *string    `gorm:"type:varchar(255)" json:"-"`
	EmailVerificationExpires *time.Time `gorm:"type:timestamptz" json:"-"`
	PasswordResetToken       *string    `gorm:"type:varchar(255)" json:"-"`
	PasswordResetExpires     *time.Time `gorm:"type:timestamptz" json:"-"`
	LastLogin                *time.Time `gorm:"type:timestamptz" json:"lastLogin,omitempty"`
	CreatedAt                time.Time  `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt                time.Time  `gorm:"type:timestamptz;not null" json:"updatedAt"`
	DeletedAt                *time.Time `gorm:"type:timestamptz;index" json:"deletedAt,omitempty"`
}

// TableName overrides default table name
func (User) TableName() string {
	return "users"
}
