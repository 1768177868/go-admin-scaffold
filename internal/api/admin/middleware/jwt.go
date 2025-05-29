package middleware

import (
	"strings"

	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWT middleware validates JWT tokens
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Get auth service
		authSvc := c.MustGet("authService").(*services.AuthService)

		// Validate token
		claims, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Get user from claims
		user, err := authSvc.GetUserFromClaims(c.Request.Context(), claims)
		if err != nil {
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Set user in context
		if user == nil {
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
