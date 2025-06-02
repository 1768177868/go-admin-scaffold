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

	// Check if user has admin role
	var count int64
	err := s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'admin' AND roles.status = 1", userModel.ID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	// Admin role has all permissions
	if count > 0 {
		return true, nil
	}

	// Check if user has the permission through their roles
	err = s.db.WithContext(ctx).Table("user_roles").
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
	// 首先检查用户是否是超级管理员
	var isAdmin int64
	err := s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'admin' AND roles.status = 1", userID).
		Count(&isAdmin).Error
	if err != nil {
		return nil, err
	}

	// 如果是超级管理员，返回所有启用的权限
	if isAdmin > 0 {
		var allPermissions []string
		err := s.db.WithContext(ctx).Model(&models.Permission{}).
			Where("status = 1").
			Pluck("name", &allPermissions).Error
		return allPermissions, err
	}

	// 如果不是超级管理员，返回分配的权限
	var permissions []string
	err = s.db.WithContext(ctx).Table("user_roles").
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

	// 首先检查用户是否有admin角色
	var hasAdminRole bool
	err := s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'admin' AND roles.status = 1", userID).
		Limit(1).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	hasAdminRole = len(roles) > 0

	if hasAdminRole {
		// 如果是admin角色，预加载所有启用的权限
		err = s.db.WithContext(ctx).
			Preload("Permissions", "status = 1").
			Joins("JOIN user_roles ON roles.id = user_roles.role_id").
			Where("user_roles.user_id = ? AND roles.status = 1", userID).
			Find(&roles).Error

		// 为admin角色添加所有启用的权限
		for i := range roles {
			if roles[i].Code == "admin" {
				var allPermissions []models.Permission
				if err := s.db.WithContext(ctx).Where("status = 1").Find(&allPermissions).Error; err != nil {
					return nil, err
				}
				roles[i].Permissions = allPermissions
			}
		}
	} else {
		// 如果不是admin角色，只加载分配的权限
		err = s.db.WithContext(ctx).
			Preload("Permissions", "status = 1").
			Joins("JOIN user_roles ON roles.id = user_roles.role_id").
			Where("user_roles.user_id = ? AND roles.status = 1", userID).
			Find(&roles).Error
	}

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
