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

// Response codes
const (
	CodeSuccess            = 0     // Success
	CodeParamError         = 10000 // Parameter error
	CodeValidationError    = 10001 // Validation error
	CodeServerError        = 10002 // Server error
	CodeNotFound           = 10003 // Not found
	CodeBusinessError      = 10004 // Business error
	CodeUnauthorized       = 10005 // Unauthorized
	CodeForbidden          = 10006 // Forbidden
	CodeCaptchaError       = 10007 // Captcha error
	CodeInvalidCaptcha     = 10008 // Invalid captcha
	CodeInvalidCredentials = 10009 // Invalid credentials
	CodeEmailTaken         = 10010 // Email already taken
	CodePermissionDenied   = 10011 // Permission denied
)

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, CodeSuccess, "success", data)
}

// Error sends an error response with trace ID
func Error(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: utils.GetTraceID(c),
	}
	c.JSON(http.StatusOK, resp)
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string) {
	Error(c, CodeValidationError, message)
}

// NotFoundError sends a not found error response
func NotFoundError(c *gin.Context) {
	Error(c, CodeNotFound, "Resource not found")
}

// BusinessError sends a business error response
func BusinessError(c *gin.Context, message string) {
	Error(c, CodeBusinessError, message)
}

// PageSuccess sends a successful paginated response
func PageSuccess(c *gin.Context, data interface{}, total int64, page int, pageSize int) {
	Success(c, gin.H{
		"items": data,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
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
	Error(c, CodeUnauthorized, "Unauthorized")
}

// ForbiddenError returns a forbidden error response
func ForbiddenError(c *gin.Context) {
	Error(c, CodeForbidden, "Forbidden")
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, CodeUnauthorized, message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, CodeForbidden, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

// ParamError sends a 400 Bad Request response for parameter errors
func ParamError(c *gin.Context, message string) {
	Error(c, CodeParamError, message)
}
