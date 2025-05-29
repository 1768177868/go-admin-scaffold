package v1

import (
	"app/internal/core/models"
	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser handles the request to get current user's profile
func GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.UnauthorizedError(c)
		return
	}

	userModel := user.(*models.User)

	// Get user with roles and permissions
	userSvc := c.MustGet("userService").(*services.UserService)
	fullUser, err := userSvc.GetByID(c.Request.Context(), userModel.ID)
	if err != nil {
		response.NotFoundError(c)
		return
	}

	// Get user permissions
	rbacSvc := c.MustGet("rbacService").(*services.RBACService)
	permissions, err := rbacSvc.GetUserPermissions(c.Request.Context(), userModel.ID)
	if err != nil {
		response.ServerError(c)
		return
	}

	result := map[string]interface{}{
		"user":        fullUser,
		"permissions": permissions,
	}

	response.Success(c, result)
}

// UpdateCurrentUser handles the request to update current user's profile
func UpdateCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.UnauthorizedError(c)
		return
	}

	userModel := user.(*models.User)

	var req struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userSvc := c.MustGet("userService").(*services.UserService)
	updateReq := &services.UpdateUserRequest{
		Email:    req.Email,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}

	updatedUser, err := userSvc.Update(c.Request.Context(), userModel.ID, updateReq)
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to update profile")
		return
	}

	response.Success(c, updatedUser)
}
