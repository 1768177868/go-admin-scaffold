package middleware

import (
	"net/http"

	"app/internal/core/models"
	"app/internal/core/repositories"
	"app/pkg/database"
	"app/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RBAC returns a middleware that checks role-based permissions
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Get user with roles
		userRepo := repositories.NewUserRepository(database.GetDB())
		user, err := userRepo.FindByID(c.Request.Context(), userID.(uint))
		if err != nil {
			logger.Error(c.Request.Context(), "Failed to get user for RBAC check", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
			c.Abort()
			return
		}

		// Get the endpoint path and method
		endpoint := c.Request.URL.Path
		method := c.Request.Method

		// Check if user has required permissions
		if !hasPermission(user, endpoint, method) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermission checks if the user has permission to access the endpoint
func hasPermission(user *models.User, endpoint, method string) bool {
	// Super admin has all permissions
	for _, role := range user.Roles {
		if role.Code == "super_admin" {
			return true
		}
	}

	// TODO: Implement more granular permission checking based on endpoint and method
	switch endpoint {
	case "/api/admin/v1/users":
		if method == "GET" {
			return hasRole(user, []string{"admin", "user_manager", "user_viewer"})
		}
		return hasRole(user, []string{"admin", "user_manager"})
	case "/api/admin/v1/roles":
		if method == "GET" {
			return hasRole(user, []string{"admin", "role_viewer"})
		}
		return hasRole(user, []string{"admin"})
	default:
		return true
	}
}

// hasRole checks if the user has any of the specified roles
func hasRole(user *models.User, roles []string) bool {
	for _, userRole := range user.Roles {
		for _, role := range roles {
			if userRole.Code == role {
				return true
			}
		}
	}
	return false
}
