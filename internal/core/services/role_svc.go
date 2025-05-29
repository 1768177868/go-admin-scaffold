package services

import (
	"app/internal/core/models"
	"app/internal/core/repositories"
	"context"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
	logSvc   *LogService
}

type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Permissions []string `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Permissions []string `json:"permissions"`
}

func NewRoleService(roleRepo *repositories.RoleRepository, logSvc *LogService) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		logSvc:   logSvc,
	}
}

func (s *RoleService) List(ctx context.Context, pagination *models.Pagination) ([]models.Role, error) {
	return s.roleRepo.List(ctx, pagination)
}

func (s *RoleService) Create(ctx context.Context, req *CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		PermList:    req.Permissions,
	}
	return s.roleRepo.Create(ctx, role)
}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(ctx, id)
}

func (s *RoleService) Update(ctx context.Context, id uint, req *UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != 0 {
		role.Status = req.Status
	}
	if len(req.Permissions) > 0 {
		role.PermList = req.Permissions
	}

	return s.roleRepo.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	return s.roleRepo.Delete(ctx, id)
}
