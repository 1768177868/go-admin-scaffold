package v1

import (
	"app/pkg/database"
	"app/pkg/redis"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	// Check database connection
	dbStatus := "healthy"
	if err := database.GetDB().Raw("SELECT 1").Error; err != nil {
		dbStatus = "unhealthy"
	}

	// Check Redis connection
	redisStatus := "healthy"
	if err := redis.GetClient().Ping(c.Request.Context()).Err(); err != nil {
		redisStatus = "unhealthy"
	}

	response.Success(c, gin.H{
		"info": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
		"status": "ok",
	})
}
