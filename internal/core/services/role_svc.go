package services

import (
	"app/internal/core/models"
	"app/internal/core/repositories"
	"context"
	"errors"

	"gorm.io/gorm"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
	db       *gorm.DB
	logSvc   *LogService
}

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	Description   string `json:"description"`
	Status        int    `json:"status"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	Status        int    `json:"status"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

func NewRoleService(roleRepo *repositories.RoleRepository, db *gorm.DB, logSvc *LogService) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		db:       db,
		logSvc:   logSvc,
	}
}

func (s *RoleService) List(ctx context.Context, pagination *models.Pagination) ([]models.Role, error) {
	return s.roleRepo.List(ctx, pagination)
}

func (s *RoleService) Create(ctx context.Context, req *CreateRoleRequest) (*models.Role, error) {
	var result *models.Role
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		role := &models.Role{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Status:      req.Status,
		}

		// Create role
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		// Assign permissions
		if len(req.PermissionIDs) > 0 {
			if err := s.assignPermissions(tx, role.ID, req.PermissionIDs); err != nil {
				return err
			}
		}

		// Load permissions
		if err := tx.Preload("Permissions").First(role, role.ID).Error; err != nil {
			return err
		}

		result = role
		return nil
	})

	return result, err
}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*models.Role, error) {
	var role models.Role
	err := s.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) Update(ctx context.Context, id uint, req *UpdateRoleRequest) (*models.Role, error) {
	var result *models.Role
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role models.Role
		if err := tx.First(&role, id).Error; err != nil {
			return err
		}

		// 禁止修改超级管理员角色
		if role.Code == "admin" {
			return errors.New("cannot modify admin role")
		}

		// Update role fields
		if req.Name != "" {
			role.Name = req.Name
		}
		if req.Code != "" {
			role.Code = req.Code
		}
		if req.Description != "" {
			role.Description = req.Description
		}
		if req.Status != 0 {
			role.Status = req.Status
		}

		// Save role
		if err := tx.Save(&role).Error; err != nil {
			return err
		}

		// Update permissions if provided
		if req.PermissionIDs != nil {
			// Remove existing permissions
			if err := tx.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
				return err
			}

			// Add new permissions
			if len(req.PermissionIDs) > 0 {
				if err := s.assignPermissions(tx, role.ID, req.PermissionIDs); err != nil {
					return err
				}
			}
		}

		// Load permissions
		if err := tx.Preload("Permissions").First(&role, role.ID).Error; err != nil {
			return err
		}

		result = &role
		return nil
	})

	return result, err
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查是否是超级管理员角色
		var role models.Role
		if err := tx.First(&role, id).Error; err != nil {
			return err
		}

		if role.Code == "admin" {
			return errors.New("cannot delete admin role")
		}

		// Remove role-permission associations
		if err := tx.Where("role_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// Remove user-role associations
		if err := tx.Exec("DELETE FROM user_roles WHERE role_id = ?", id).Error; err != nil {
			return err
		}

		// Delete role
		return tx.Delete(&models.Role{}, id).Error
	})
}

// assignPermissions assigns permissions to a role
func (s *RoleService) assignPermissions(tx *gorm.DB, roleID uint, permissionIDs []uint) error {
	var rolePermissions []models.RolePermission
	for _, permID := range permissionIDs {
		rolePermissions = append(rolePermissions, models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		})
	}
	return tx.Create(&rolePermissions).Error
}

// GetPermissions returns all permissions for a role
func (s *RoleService) GetPermissions(ctx context.Context, roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.status = 1", roleID).
		Find(&permissions).Error
	return permissions, err
}

// UpdatePermissions updates the permissions of a role
func (s *RoleService) UpdatePermissions(ctx context.Context, roleID uint, req *UpdateRolePermissionsRequest) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if role exists
		var role models.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		// 禁止修改超级管理员角色权限
		if role.Code == "admin" {
			return errors.New("cannot modify admin role permissions")
		}

		// Filter out invalid permission IDs (0 or duplicates)
		validPermIDs := make([]uint, 0)
		seen := make(map[uint]bool)
		for _, id := range req.PermissionIDs {
			if id > 0 && !seen[id] {
				// Check if permission exists
				var count int64
				if err := tx.Model(&models.Permission{}).Where("id = ? AND status = 1", id).Count(&count).Error; err != nil {
					return err
				}
				if count > 0 {
					validPermIDs = append(validPermIDs, id)
					seen[id] = true
				}
			}
		}

		// Remove existing permissions
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// Add new permissions
		if len(validPermIDs) > 0 {
			return s.assignPermissions(tx, roleID, validPermIDs)
		}

		return nil
	})
}
