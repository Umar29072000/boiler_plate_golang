package repository

import (
	"context"

	"boiler_plate_be_golang/app/config"
	model "boiler_plate_be_golang/pkg/model/database"

	"gorm.io/gorm"
)

// IUserRepository defines user repository interface
type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
	FindByVerificationToken(ctx context.Context, token string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
}

// UserRepository handles user data access
type UserRepository struct {
	db        *gorm.DB
	AppConfig config.App
}

// NewUserRepository creates new user repository
func NewUserRepository(db *gorm.DB, appConfig config.App) *UserRepository {
	return &UserRepository{
		db:        db,
		AppConfig: appConfig,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll finds all users with pagination
func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Update updates user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes user (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// Count counts total users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

// FindByVerificationToken finds user by email verification token
func (r *UserRepository) FindByVerificationToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email_verification_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByResetToken finds user by password reset token
func (r *UserRepository) FindByResetToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("password_reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
