package v1

import (
	"net/http"
	"strconv"

	"app/internal/core/models"
	"app/internal/core/services"

	"github.com/gin-gonic/gin"
)

// ListLoginLogs handles the request to get a paginated list of login logs
func ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse query parameters
	query := &services.LogQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logSvc := c.MustGet("logService").(*services.LogService)
	logs, err := logSvc.ListLoginLogs(c.Request.Context(), pagination, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch login logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
			"total":     pagination.Total,
		},
	})
}

// ListOperationLogs handles the request to get a paginated list of operation logs
func ListOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse query parameters
	query := &services.LogQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logSvc := c.MustGet("logService").(*services.LogService)
	logs, err := logSvc.ListOperationLogs(c.Request.Context(), pagination, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch operation logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
			"total":     pagination.Total,
		},
	})
}

// GetUserLogs handles the request to get a user's login and operation history
func GetUserLogs(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	logSvc := c.MustGet("logService").(*services.LogService)

	loginLogs, err := logSvc.GetUserLoginHistory(c.Request.Context(), uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch login history"})
		return
	}

	operationLogs, err := logSvc.GetUserOperationHistory(c.Request.Context(), uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch operation history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login_logs":     loginLogs,
		"operation_logs": operationLogs,
	})
}
