package utils

import "github.com/gin-gonic/gin"

// GetTraceID returns the trace ID from the context
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}
