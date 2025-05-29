package v1

import (
	"app/internal/core/services"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// Login handles user authentication and returns a JWT token
func Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	authSvc := c.MustGet("authService").(*services.AuthService)
	resp, err := authSvc.Login(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			response.Error(c, response.CodeInvalidCredentials, "invalid credentials")
			return
		}
		if err == services.ErrUserInactive {
			response.Error(c, response.CodeForbidden, "user is inactive")
			return
		}
		response.ServerError(c)
		return
	}

	response.Success(c, resp)
}

// RefreshToken handles token refresh requests
func RefreshToken(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedError(c)
		return
	}

	authSvc := c.MustGet("authService").(*services.AuthService)
	token, err := authSvc.RefreshToken(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, response.CodeServerError, "failed to refresh token")
		return
	}

	response.Success(c, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	})
}
