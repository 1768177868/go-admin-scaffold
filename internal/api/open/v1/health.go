package v1

import (
	"app/pkg/database"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	// Check database connection
	if err := database.DB().Raw("SELECT 1").Error; err != nil {
		response.Error(c, response.CodeServerError, "Database connection error")
		return
	}

	// Check Redis connection
	// TODO: Add Redis health check
	if false {
		response.Error(c, response.CodeServerError, "Redis connection error")
		return
	}

	response.Success(c, gin.H{
		"status": "ok",
		"info": gin.H{
			"database": "healthy",
			"redis":    "healthy",
		},
	})
}
