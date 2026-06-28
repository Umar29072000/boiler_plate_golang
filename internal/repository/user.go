package repository

import (
	"context"
	"errors"
	"math"

	"boiler_plate_be_golang/app/config"
	"boiler_plate_be_golang/domains/dto"
	"boiler_plate_be_golang/internal/rest/response"
	model "boiler_plate_be_golang/pkg/model/database"

	"gorm.io/gorm"
)

// IUserRepository defines user repository interface
type IUserRepository interface {
	Show(ctx context.Context, req dto.GetUserData) (res response.PaginationResponse, err error)
	Create(ctx context.Context, req dto.CreateUserData) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, req dto.UpdateUserData) (*model.User, error)
	Delete(ctx context.Context, id string) error
	FindByVerificationToken(ctx context.Context, token string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
}

// UserRepository handles user data access
type UserRepository struct {
	DB        *gorm.DB
	AppConfig config.App
}

// NewUserRepository creates new user repository
func NewUserRepository(db *gorm.DB, appConfig config.App) *UserRepository {
	return &UserRepository{
		DB:        db,
		AppConfig: appConfig,
	}
}

// Show gets users with pagination
func (r *UserRepository) Show(ctx context.Context, req dto.GetUserData) (res response.PaginationResponse, err error) {
	var users []model.User
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.User{})

	// Apply filters
	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}
	if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	}
	if req.Search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Count total items if not disabled
	if req.DisableCalculateTotal != "true" {
		if err = query.Count(&total).Error; err != nil {
			return res, err
		}
	}

	// Apply sorting
	if req.Field != "" && req.Sort != "" {
		query = query.Order(req.Field + " " + req.Sort)
	} else {
		query = query.Order("created_at desc")
	}

	// Apply pagination
	if req.Page > 0 && req.Limit > 0 {
		offset := (req.Page - 1) * req.Limit
		query = query.Offset(offset).Limit(req.Limit)
	}

	if err = query.Find(&users).Error; err != nil {
		return res, err
	}

	// Calculate total pages
	var totalPages int64
	if req.Limit > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(req.Limit)))
	}

	res = response.PaginationResponse{
		Items:      users,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       int64(req.Page),
		Limit:      int64(req.Limit),
	}

	return res, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, req dto.CreateUserData) (*model.User, error) {
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates user
func (r *UserRepository) Update(ctx context.Context, req dto.UpdateUserData) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("id = ?", req.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("USER_NOT_FOUND")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := r.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete deletes user (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.DB.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("USER_NOT_FOUND")
	}
	return nil
}

// FindByVerificationToken finds user by email verification token
func (r *UserRepository) FindByVerificationToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("email_verification_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByResetToken finds user by password reset token
func (r *UserRepository) FindByResetToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("password_reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
