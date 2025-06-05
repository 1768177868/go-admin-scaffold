package services

import (
	"app/internal/core/models"
	"context"

	"gorm.io/gorm"
)

// RBACService handles role-based access control
type RBACService struct {
	db      *gorm.DB
	authSvc AuthServiceInterface
}

// NewRBACService creates a new RBAC service instance
func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{
		db: db,
	}
}

// SetAuthService sets the auth service instance
func (s *RBACService) SetAuthService(authSvc AuthServiceInterface) {
	s.authSvc = authSvc
}

// CheckPermission checks if a user has the specified permission
func (s *RBACService) CheckPermission(ctx context.Context, user interface{}, permission string) (bool, error) {
	// Type assertion to get user model
	userModel, ok := user.(*models.User)
	if !ok {
		return false, nil
	}

	// Check if user is super admin using auth service
	if s.authSvc != nil && s.authSvc.IsSuperAdmin(userModel.ID) {
		return true, nil
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

	// Check if user has the permission through their role-menu associations
	err = s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN role_menus ON user_roles.role_id = role_menus.role_id").
		Joins("JOIN menus ON role_menus.menu_id = menus.id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND menus.permission = ? AND menus.status = 1 AND menus.visible = 1 AND roles.status = 1",
			userModel.ID, permission).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserPermissions returns all permissions for a user
func (s *RBACService) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	// Check if user is super admin using auth service
	if s.authSvc != nil && s.authSvc.IsSuperAdmin(userID) {
		var allPermissions []string
		err := s.db.WithContext(ctx).Model(&models.Menu{}).
			Where("status = 1 AND visible = 1 AND permission != ''").
			Pluck("permission", &allPermissions).Error
		return allPermissions, err
	}

	// 检查用户是否是管理员
	var isAdmin int64
	err := s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'admin' AND roles.status = 1", userID).
		Count(&isAdmin).Error
	if err != nil {
		return nil, err
	}

	// 如果是管理员，返回所有启用的权限
	if isAdmin > 0 {
		var allPermissions []string
		err := s.db.WithContext(ctx).Model(&models.Menu{}).
			Where("status = 1 AND visible = 1 AND permission != ''").
			Pluck("permission", &allPermissions).Error
		return allPermissions, err
	}

	// 如果不是超级管理员或管理员，返回分配的权限
	var permissions []string
	err = s.db.WithContext(ctx).Table("user_roles").
		Select("DISTINCT menus.permission").
		Joins("JOIN role_menus ON user_roles.role_id = role_menus.role_id").
		Joins("JOIN menus ON role_menus.menu_id = menus.id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND menus.status = 1 AND menus.visible = 1 AND menus.permission != '' AND roles.status = 1", userID).
		Pluck("menus.permission", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetUserRoles returns all roles for a user with their menus
func (s *RBACService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	var roles []models.Role

	// Check if user is super admin using auth service
	isSuperAdmin := s.authSvc != nil && s.authSvc.IsSuperAdmin(userID)

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

	if hasAdminRole || isSuperAdmin {
		// 如果是admin角色或超级管理员，预加载所有启用的菜单
		err = s.db.WithContext(ctx).
			Preload("Menus", "status = 1 AND visible = 1").
			Joins("JOIN user_roles ON roles.id = user_roles.role_id").
			Where("user_roles.user_id = ? AND roles.status = 1", userID).
			Find(&roles).Error

		// 为admin角色或超级管理员添加所有启用的菜单
		for i := range roles {
			if roles[i].Code == "admin" || isSuperAdmin {
				var allMenus []models.Menu
				if err := s.db.WithContext(ctx).Where("status = 1 AND visible = 1").Find(&allMenus).Error; err != nil {
					return nil, err
				}
				roles[i].Menus = allMenus
			}
		}
	} else {
		// 如果不是admin角色或超级管理员，只加载分配的菜单
		err = s.db.WithContext(ctx).
			Preload("Menus", "status = 1 AND visible = 1").
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
func (s *RBACService) HasAnyPermission(ctx context.Context, userID uint, permissions []string) (bool, error) {
	if len(permissions) == 0 {
		return false, nil
	}

	// Check if user is super admin
	if s.authSvc != nil && s.authSvc.IsSuperAdmin(userID) {
		return true, nil
	}

	// Check if user has admin role
	var count int64
	err := s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'admin' AND roles.status = 1", userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	// Check if user has any of the specified permissions through role-menu associations
	err = s.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN role_menus ON user_roles.role_id = role_menus.role_id").
		Joins("JOIN menus ON role_menus.menu_id = menus.id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND menus.permission IN ? AND menus.status = 1 AND menus.visible = 1 AND roles.status = 1",
			userID, permissions).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
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
