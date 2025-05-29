package v1

import (
	"strconv"

	"app/internal/core/models"
	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
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
