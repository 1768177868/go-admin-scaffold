package v1

import (
	"strconv"

	"app/internal/core/models"
	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListUsers handles the request to get a paginated list of users
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	users, err := userSvc.List(c.Request.Context(), pagination)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to fetch users")
		return
	}

	response.PageSuccess(c, users, pagination.Total, pagination.Page, pagination.PageSize)
}

// CreateUser handles the request to create a new user
func CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	user, err := userSvc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to create user")
		return
	}

	response.Success(c, user)
}

// GetUser handles the request to get a user by ID
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid user ID")
		return
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	user, err := userSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFoundError(c)
		return
	}

	response.Success(c, user)
}

// UpdateUser handles the request to update a user
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid user ID")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	user, err := userSvc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to update user")
		return
	}

	response.Success(c, user)
}

// DeleteUser handles the request to delete a user
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid user ID")
		return
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	if err := userSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, response.CodeServerError, "failed to delete user")
		return
	}

	response.Success(c, nil)
}

// ExportUsers handles the request to export user list data
func ExportUsers(c *gin.Context) {
	var req services.ExportUserListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	users, err := userSvc.ExportUserList(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to export users")
		return
	}

	response.Success(c, users)
}
