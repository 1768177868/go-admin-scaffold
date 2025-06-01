package services

import (
	"app/internal/core/models"
	"context"

	"gorm.io/gorm"
)

type PermissionService struct {
	db *gorm.DB
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Module      string `json:"module" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Resource    string `json:"resource" binding:"required"`
	Status      int    `json:"status"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Module      string `json:"module"`
	Action      string `json:"action"`
	Resource    string `json:"resource"`
	Status      int    `json:"status"`
}

// PermissionTreeNode represents a node in the permission tree
type PermissionTreeNode struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	DisplayName string               `json:"display_name"`
	Description string               `json:"description"`
	Module      string               `json:"module"`
	Action      string               `json:"action"`
	Resource    string               `json:"resource"`
	Status      int                  `json:"status"`
	Children    []PermissionTreeNode `json:"children,omitempty"`
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{
		db: db,
	}
}

func (s *PermissionService) List(ctx context.Context, pagination *models.Pagination) ([]models.Permission, error) {
	var permissions []models.Permission
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Permission{})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.PageSize
		query = query.Offset(offset).Limit(pagination.PageSize)
		pagination.Total = total
	}

	// Get permissions
	err := query.Find(&permissions).Error
	return permissions, err
}

func (s *PermissionService) Create(ctx context.Context, req *CreatePermissionRequest) (*models.Permission, error) {
	permission := &models.Permission{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Module:      req.Module,
		Action:      req.Action,
		Resource:    req.Resource,
		Status:      req.Status,
	}

	err := s.db.WithContext(ctx).Create(permission).Error
	return permission, err
}

func (s *PermissionService) GetByID(ctx context.Context, id uint) (*models.Permission, error) {
	var permission models.Permission
	err := s.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (s *PermissionService) Update(ctx context.Context, id uint, req *UpdatePermissionRequest) (*models.Permission, error) {
	var permission models.Permission
	if err := s.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.DisplayName != "" {
		permission.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	if req.Module != "" {
		permission.Module = req.Module
	}
	if req.Action != "" {
		permission.Action = req.Action
	}
	if req.Resource != "" {
		permission.Resource = req.Resource
	}
	if req.Status != 0 {
		permission.Status = req.Status
	}

	err := s.db.WithContext(ctx).Save(&permission).Error
	return &permission, err
}

func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove role-permission associations
		if err := tx.Where("permission_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// Delete permission
		return tx.Delete(&models.Permission{}, id).Error
	})
}

// GetByModule returns permissions grouped by module
func (s *PermissionService) GetByModule(ctx context.Context) (map[string][]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.WithContext(ctx).Where("status = 1").Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]models.Permission)
	for _, perm := range permissions {
		result[perm.Module] = append(result[perm.Module], perm)
	}

	return result, nil
}

// GetByNames returns permissions by their names
func (s *PermissionService) GetByNames(ctx context.Context, names []string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.WithContext(ctx).Where("name IN ? AND status = 1", names).Find(&permissions).Error
	return permissions, err
}

// GetPermissionTree returns the permission tree structure
func (s *PermissionService) GetPermissionTree(ctx context.Context) ([]PermissionTreeNode, error) {
	var permissions []models.Permission
	err := s.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	// First, group by module
	moduleMap := make(map[string][]models.Permission)
	for _, perm := range permissions {
		moduleMap[perm.Module] = append(moduleMap[perm.Module], perm)
	}

	// Build tree
	var tree []PermissionTreeNode
	for module, perms := range moduleMap {
		// Create module node
		moduleNode := PermissionTreeNode{
			Name:        module,
			DisplayName: module, // You might want to add a display name mapping
			Module:      module,
			Children:    make([]PermissionTreeNode, 0),
		}

		// Group by resource
		resourceMap := make(map[string][]models.Permission)
		for _, perm := range perms {
			resourceMap[perm.Resource] = append(resourceMap[perm.Resource], perm)
		}

		// Add resource nodes
		for resource, resourcePerms := range resourceMap {
			resourceNode := PermissionTreeNode{
				Name:        resource,
				DisplayName: resource, // You might want to add a display name mapping
				Resource:    resource,
				Children:    make([]PermissionTreeNode, 0),
			}

			// Add permission nodes
			for _, perm := range resourcePerms {
				permNode := PermissionTreeNode{
					ID:          perm.ID,
					Name:        perm.Name,
					DisplayName: perm.DisplayName,
					Description: perm.Description,
					Module:      perm.Module,
					Action:      perm.Action,
					Resource:    perm.Resource,
					Status:      perm.Status,
				}
				resourceNode.Children = append(resourceNode.Children, permNode)
			}

			moduleNode.Children = append(moduleNode.Children, resourceNode)
		}

		tree = append(tree, moduleNode)
	}

	return tree, nil
}
