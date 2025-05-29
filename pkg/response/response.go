package response

import (
	"net/http"

	"app/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Response represents the unified response structure
type Response struct {
	Code    int         `json:"code"`     // Business status code
	Message string      `json:"message"`  // Response message
	Data    interface{} `json:"data"`     // Response data
	TraceID string      `json:"trace_id"` // Trace ID for request tracking
}

// PageData represents paginated data
type PageData struct {
	List     interface{} `json:"list"`      // List data
	Total    int64       `json:"total"`     // Total count
	Page     int         `json:"page"`      // Current page
	PageSize int         `json:"page_size"` // Page size
	Pages    int         `json:"pages"`     // Total pages
}

const (
	// Success code
	CodeSuccess = 0

	// Common error codes (10000-10999)
	CodeServerError     = 10000 // Server error
	CodeParamError      = 10001 // Parameter error
	CodeValidationError = 10002 // Validation error
	CodeNotFound        = 10003 // Not found
	CodeBusinessError   = 10004 // Business error
)

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, CodeSuccess, "success", data)
}

// Error sends an error response
func Error(c *gin.Context, code int, message string) {
	JSON(c, http.StatusOK, code, message, nil)
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string) {
	JSON(c, http.StatusOK, CodeValidationError, message, nil)
}

// NotFoundError sends a not found error response
func NotFoundError(c *gin.Context) {
	JSON(c, http.StatusOK, CodeNotFound, "Resource not found", nil)
}

// BusinessError sends a business error response
func BusinessError(c *gin.Context, message string) {
	JSON(c, http.StatusOK, CodeBusinessError, message, nil)
}

// PageSuccess sends a successful paginated response
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	pageData := &PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}

	Success(c, pageData)
}

// JSON sends a JSON response with trace ID
func JSON(c *gin.Context, httpStatus, code int, message string, data interface{}) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
		TraceID: utils.GetTraceID(c),
	}
	c.JSON(httpStatus, resp)
}

// Page represents paginated data
type Page struct {
	List     interface{} `json:"list"`      // List data
	Total    int64       `json:"total"`     // Total count
	Page     int         `json:"page"`      // Current page
	PageSize int         `json:"page_size"` // Page size
	Pages    int         `json:"pages"`     // Total pages
}

// ServerError returns a server error response
func ServerError(c *gin.Context) {
	Error(c, CodeServerError, "Internal server error")
}

// UnauthorizedError returns an unauthorized error response
func UnauthorizedError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: "Unauthorized",
		Data:    nil,
	})
}

// ForbiddenError returns a forbidden error response
func ForbiddenError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: "Forbidden",
		Data:    nil,
	})
}

// ParamError returns a parameter error response
func ParamError(c *gin.Context, message string) {
	Error(c, CodeParamError, message)
}
