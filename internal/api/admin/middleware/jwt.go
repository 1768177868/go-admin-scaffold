package middleware

import (
	"log"
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
			log.Printf("[ERROR] Invalid authorization header format: %s", authHeader)
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Get auth service and user service
		authSvc := c.MustGet("authService").(*services.AuthService)

		tokenString := parts[1]
		log.Printf("[DEBUG] Received JWT token: %s", tokenString)
		log.Printf("[DEBUG] JWT token length: %d", len(tokenString))
		log.Printf("[DEBUG] JWT token segments: %d", len(strings.Split(tokenString, ".")))

		// Validate token
		claims, err := authSvc.ValidateToken(tokenString)
		if err != nil {
			log.Printf("[ERROR] Failed to validate token: %v", err)
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Get user from claims
		user, err := authSvc.GetUserFromClaims(c.Request.Context(), claims)
		if err != nil {
			log.Printf("[ERROR] Failed to get user from claims: %v", err)
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Set user in context
		if user == nil {
			log.Printf("[ERROR] User is nil after GetUserFromClaims")
			response.UnauthorizedError(c)
			c.Abort()
			return
		}

		// Set IsSuperAdmin field
		user.IsSuperAdmin = authSvc.IsSuperAdmin(user.ID)
		log.Printf("[DEBUG] Set IsSuperAdmin field for user %d: %v", user.ID, user.IsSuperAdmin)

		c.Set("user", user)
		c.Next()
	}
}
