package services

import (
	"app/internal/core/models"
	"context"

	"gorm.io/gorm"
)

// RBACService handles role-based access control
type RBACService struct {
	db *gorm.DB
}

// NewRBACService creates a new RBAC service instance
func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{
		db: db,
	}
}

// CheckPermission checks if a user has the specified permission
func (s *RBACService) CheckPermission(ctx context.Context, user interface{}, permission string) (bool, error) {
	// Type assertion to get user model
	userModel, ok := user.(*models.User)
	if !ok {
		return false, nil
	}

	// Super admin check (user ID = 1 typically has all permissions)
	if userModel.ID == 1 {
		return true, nil
	}

	// Check if user has the permission through their roles
	var count int64
	err := s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND permissions.name = ? AND permissions.status = 1 AND roles.status = 1",
			userModel.ID, permission).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserPermissions returns all permissions for a user
func (s *RBACService) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	var permissions []string

	err := s.db.WithContext(ctx).Table("user_roles").
		Select("DISTINCT permissions.name").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND permissions.status = 1 AND roles.status = 1", userID).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetUserRoles returns all roles for a user with their permissions
func (s *RBACService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	var roles []models.Role

	err := s.db.WithContext(ctx).
		Preload("Permissions", "status = 1").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = 1", userID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// HasAnyPermission checks if user has any of the specified permissions
func (s *RBACService) HasAnyPermission(ctx context.Context, user interface{}, permissions []string) (bool, error) {
	for _, perm := range permissions {
		has, err := s.CheckPermission(ctx, user, perm)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

// HasAllPermissions checks if user has all of the specified permissions
func (s *RBACService) HasAllPermissions(ctx context.Context, user interface{}, permissions []string) (bool, error) {
	for _, perm := range permissions {
		has, err := s.CheckPermission(ctx, user, perm)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}
	}
	return true, nil
}
