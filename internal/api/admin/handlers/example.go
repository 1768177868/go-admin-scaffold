package handlers

import (
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExampleRequest represents an example request
type ExampleRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// ExampleResponse represents an example response
type ExampleResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ExampleHandler handles example requests
type ExampleHandler struct {
	// Add any dependencies here
}

// NewExampleHandler creates a new example handler
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// Create handles the creation of an example resource
func (h *ExampleHandler) Create(c *gin.Context) {
	var req ExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// Check if email exists (example of business logic error)
	if req.Email == "exists@example.com" {
		response.Error(c, response.CodeEmailTaken, "Email is already taken")
		return
	}

	// Example of successful response
	resp := &ExampleResponse{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
	}
	response.Success(c, resp)
}

// List handles listing example resources with pagination
func (h *ExampleHandler) List(c *gin.Context) {
	// Example pagination parameters
	page := 1
	pageSize := 10
	total := int64(100)

	// Example list data
	list := []ExampleResponse{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
	}

	response.PageSuccess(c, list, total, page, pageSize)
}

// Get handles getting a single example resource
func (h *ExampleHandler) Get(c *gin.Context) {
	// Example of not found error
	id := c.Param("id")
	if id == "0" {
		response.NotFoundError(c)
		return
	}

	// Example of successful response
	resp := &ExampleResponse{
		ID:    1,
		Name:  "Example User",
		Email: "user@example.com",
	}
	response.Success(c, resp)
}

// Update handles updating an example resource
func (h *ExampleHandler) Update(c *gin.Context) {
	var req ExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// Example of permission error
	if req.Name == "admin" {
		response.Error(c, response.CodePermissionDenied, "Cannot modify admin user")
		return
	}

	// Example of successful response
	resp := &ExampleResponse{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
	}
	response.Success(c, resp)
}

// Delete handles deleting an example resource
func (h *ExampleHandler) Delete(c *gin.Context) {
	// Example of business logic error
	id := c.Param("id")
	if id == "1" {
		response.BusinessError(c, "Cannot delete the default user")
		return
	}

	response.Success(c, nil)
}

// HandleError demonstrates different error responses
func (h *ExampleHandler) HandleError(c *gin.Context) {
	errorType := c.Query("type")

	switch errorType {
	case "validation":
		response.ValidationError(c, "Invalid input data")
	case "unauthorized":
		response.UnauthorizedError(c)
	case "forbidden":
		response.ForbiddenError(c)
	case "notfound":
		response.NotFoundError(c)
	case "param":
		response.ParamError(c, "Missing required parameter")
	case "business":
		response.BusinessError(c, "Business rule violation")
	case "server":
		response.ServerError(c)
	default:
		response.Success(c, map[string]string{"message": "No error"})
	}
}
