package response

// Common business status codes
const (
	// Common error codes (10xxx)
	CodeUnauthorized = 10401
	CodeForbidden    = 10403

	// Authentication errors (11xxx)
	CodeInvalidCredentials = 11001

	// User errors (12xxx)
	CodeEmailTaken       = 12008
	CodePermissionDenied = 13003
)

// CodeMessages maps error codes to their default messages
var CodeMessages = map[int]string{
	CodeUnauthorized:       "Unauthorized",
	CodeForbidden:          "Forbidden",
	CodeInvalidCredentials: "Invalid credentials",
	CodeEmailTaken:         "Email is already taken",
	CodePermissionDenied:   "Permission denied",
}
