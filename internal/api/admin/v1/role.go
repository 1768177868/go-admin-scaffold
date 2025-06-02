package v1

import (
	"strconv"

	"app/internal/core/models"
	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListRoles handles the request to get a paginated list of roles
func ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	roles, err := roleSvc.List(c.Request.Context(), pagination)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to fetch roles")
		return
	}

	response.PageSuccess(c, roles, pagination.Total, pagination.Page, pagination.PageSize)
}

// CreateRole handles the request to create a new role
func CreateRole(c *gin.Context) {
	var req services.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	role, err := roleSvc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to create role")
		return
	}

	response.Success(c, role)
}

// GetRole handles the request to get a role by ID
func GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid role ID")
		return
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	role, err := roleSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFoundError(c)
		return
	}

	response.Success(c, role)
}

// UpdateRole handles the request to update a role
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid role ID")
		return
	}

	var req services.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	role, err := roleSvc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to update role")
		return
	}

	response.Success(c, role)
}

// DeleteRole handles the request to delete a role
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid role ID")
		return
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	if err := roleSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, response.CodeServerError, "failed to delete role")
		return
	}

	response.Success(c, nil)
}

// RolePermissionResponse represents the response for role permissions
type RolePermissionResponse struct {
	PermissionTree []PermissionTreeNode `json:"permission_tree"`
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
	Assigned    bool                 `json:"assigned"`
	Children    []PermissionTreeNode `json:"children,omitempty"`
}

// GetRolePermissions handles the request to get permissions of a role
func GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid role ID")
		return
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	permissionSvc := c.MustGet("permissionService").(*services.PermissionService)

	// Get all permissions
	allPermissions, err := permissionSvc.List(c.Request.Context(), &models.Pagination{
		Page:     1,
		PageSize: 1000, // 设置一个较大的数，获取所有权限
	})
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to get permissions")
		return
	}

	// Get role's permissions
	rolePermissions, err := roleSvc.GetPermissions(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to get role permissions")
		return
	}

	// Create a map for quick lookup of assigned permissions
	assignedMap := make(map[uint]bool)
	for _, p := range rolePermissions {
		assignedMap[p.ID] = true
	}

	// Group permissions by module
	moduleMap := make(map[string][]PermissionTreeNode)
	for _, p := range allPermissions {
		// Skip disabled permissions
		if p.Status != 1 {
			continue
		}
		node := PermissionTreeNode{
			ID:          p.ID,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			Module:      p.Module,
			Action:      p.Action,
			Resource:    p.Resource,
			Status:      p.Status,
			Assigned:    assignedMap[p.ID],
		}
		moduleMap[p.Module] = append(moduleMap[p.Module], node)
	}

	// Create tree structure
	var result RolePermissionResponse
	result.PermissionTree = make([]PermissionTreeNode, 0)

	// Create module nodes
	for module, permissions := range moduleMap {
		// Calculate if module should be assigned based on children
		moduleAssigned := true
		for _, p := range permissions {
			if !p.Assigned {
				moduleAssigned = false
				break
			}
		}

		moduleNode := PermissionTreeNode{
			Name:        module,
			DisplayName: module, // You might want to map this to a proper display name
			Module:      module,
			Assigned:    moduleAssigned,
			Children:    permissions,
		}
		result.PermissionTree = append(result.PermissionTree, moduleNode)
	}

	response.Success(c, result)
}

// UpdateRolePermissions handles the request to update permissions of a role
func UpdateRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid role ID")
		return
	}

	var req services.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// Filter out invalid permission IDs
	validIDs := make([]uint, 0)
	for _, pid := range req.PermissionIDs {
		if pid > 0 {
			validIDs = append(validIDs, pid)
		}
	}

	// Update request with valid IDs
	req.PermissionIDs = validIDs

	// Validate that we have at least one valid permission ID
	if len(req.PermissionIDs) == 0 {
		response.ValidationError(c, "no valid permission IDs provided")
		return
	}

	roleSvc := c.MustGet("roleService").(*services.RoleService)
	if err := roleSvc.UpdatePermissions(c.Request.Context(), uint(id), &req); err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFoundError(c)
			return
		}
		response.Error(c, response.CodeServerError, "failed to update role permissions: "+err.Error())
		return
	}

	response.Success(c, nil)
}
