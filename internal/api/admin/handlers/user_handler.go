package handlers

import (
	"app/internal/core/services"
	"app/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *services.UserService
}

func NewUserHandler(userSvc *services.UserService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

// UpdateUserRoles handles updating a user's roles
// @Summary Update user roles
// @Description Updates the roles assigned to a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param roles body UpdateUserRolesRequest true "Role IDs"
// @Success 200 {object} response.Response
// @Router /admin/v1/users/{id}/roles [put]
func (h *UserHandler) UpdateUserRoles(c *gin.Context) {
	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "Invalid user ID")
		return
	}

	// Update user roles
	if err := h.userSvc.UpdateUserRoles(c.Request.Context(), uint(userID), req.RoleIDs); err != nil {
		response.BusinessError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "User roles updated successfully"})
}
