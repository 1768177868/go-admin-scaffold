package v1

import (
	"net/http"

	"app/pkg/database"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles the health check request
func HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := database.GetDB().DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "database connection error",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "service is healthy",
	})
}
