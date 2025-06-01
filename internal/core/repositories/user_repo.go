package repositories

import (
	"context"

	"app/internal/core/models"
	"app/internal/core/types"

	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListWithRoles retrieves a paginated list of users with their roles
func (r *UserRepository) ListWithRoles(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Get total count
	if err := r.db.Model(&models.User{}).Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("last_login_at", gorm.Expr("NOW()")).
		Error
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetDB returns the database instance
func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// ListWithFilters retrieves a paginated list of users with search filters
func (r *UserRepository) ListWithFilters(ctx context.Context, pagination *models.Pagination, filters *types.UserSearchFilters) ([]models.User, error) {
	var users []models.User
	query := r.db.WithContext(ctx).Preload("Roles")

	// Apply filters
	if filters != nil {
		if filters.Username != "" {
			query = query.Where("username LIKE ?", "%"+filters.Username+"%")
		}
		if filters.Email != "" {
			query = query.Where("email LIKE ?", "%"+filters.Email+"%")
		}
		if filters.Status != nil {
			query = query.Where("status = ?", *filters.Status)
		}
		if filters.RoleID > 0 {
			query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
				Where("user_roles.role_id = ?", filters.RoleID)
		}
	}

	// Get total count for pagination
	var total int64
	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total

	// Apply pagination
	err := query.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
