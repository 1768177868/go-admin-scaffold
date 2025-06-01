package v1

import (
	"strconv"

	"app/internal/core/models"
	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListPermissions handles the request to get a paginated list of permissions
func ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	permissions, err := permSvc.List(c.Request.Context(), pagination)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to fetch permissions")
		return
	}

	response.PageSuccess(c, permissions, pagination.Total, pagination.Page, pagination.PageSize)
}

// CreatePermission handles the request to create a new permission
func CreatePermission(c *gin.Context) {
	var req services.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	permission, err := permSvc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to create permission")
		return
	}

	response.Success(c, permission)
}

// GetPermission handles the request to get a permission by ID
func GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid permission ID")
		return
	}

	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	permission, err := permSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFoundError(c)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission handles the request to update a permission
func UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid permission ID")
		return
	}

	var req services.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	permission, err := permSvc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to update permission")
		return
	}

	response.Success(c, permission)
}

// DeletePermission handles the request to delete a permission
func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid permission ID")
		return
	}

	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	if err := permSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, response.CodeServerError, "failed to delete permission")
		return
	}

	response.Success(c, nil)
}

// GetPermissionsByModule handles the request to get permissions grouped by module
func GetPermissionsByModule(c *gin.Context) {
	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	permissions, err := permSvc.GetByModule(c.Request.Context())
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to fetch permissions")
		return
	}

	response.Success(c, permissions)
}

// GetPermissionTree handles the request to get the permission tree
func GetPermissionTree(c *gin.Context) {
	permSvc := c.MustGet("permissionService").(*services.PermissionService)
	tree, err := permSvc.GetPermissionTree(c.Request.Context())
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to fetch permission tree")
		return
	}

	response.Success(c, tree)
}
