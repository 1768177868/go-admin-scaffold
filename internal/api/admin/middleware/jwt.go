package middleware

import (
	"strings"

	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	CodeUnauthorized = 10401
)

// JWT middleware validates JWT tokens
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, CodeUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, CodeUnauthorized, "Invalid authorization header")
			c.Abort()
			return
		}

		// Get auth service
		authSvc := c.MustGet("authService").(*services.AuthService)

		// Validate token
		claims, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, CodeUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Get user from claims
		user, err := authSvc.GetUserFromClaims(c.Request.Context(), claims)
		if err != nil {
			response.Error(c, CodeUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		// Set user in context
		if user == nil {
			response.Error(c, CodeUnauthorized, "User not found")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
