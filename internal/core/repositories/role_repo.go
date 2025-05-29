package repositories

import (
	"app/internal/core/models"
	"context"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) List(ctx context.Context, pagination *models.Pagination) ([]models.Role, error) {
	var roles []models.Role
	query := r.db.WithContext(ctx).Model(&models.Role{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total

	result := query.Offset((pagination.Page - 1) * pagination.PageSize).
		Limit(pagination.PageSize).
		Find(&roles)

	if result.Error != nil {
		return nil, result.Error
	}

	return roles, nil
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *models.Role) (*models.Role, error) {
	if err := r.db.WithContext(ctx).Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}
