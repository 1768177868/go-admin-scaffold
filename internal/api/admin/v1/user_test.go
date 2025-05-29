package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/internal/core/models"
	"app/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, req *services.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, id uint, req *services.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) List(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) ExportUserList(ctx context.Context, req *services.ExportUserListRequest) ([]models.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.User), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockUserService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockSvc := new(MockUserService)
	r.Use(func(c *gin.Context) {
		c.Set("userService", mockSvc)
	})
	return r, mockSvc
}

func TestListUsers(t *testing.T) {
	r, mockSvc := setupTestRouter()
	r.GET("/users", ListUsers)

	// Test case 1: Successful list
	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: 1},
			Username:  "user1",
			Email:     "user1@example.com",
		},
		{
			BaseModel: models.BaseModel{ID: 2},
			Username:  "user2",
			Email:     "user2@example.com",
		},
	}

	mockSvc.On("List", mock.Anything, mock.AnythingOfType("*models.Pagination")).
		Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["pagination"])
}

func TestExportUsers(t *testing.T) {
	r, mockSvc := setupTestRouter()
	r.POST("/users/export", ExportUsers)

	// Test case 1: Successful export
	now := time.Now()
	req := services.ExportUserListRequest{
		Username:  "user",
		Email:     "@example.com",
		Status:    &[]int{1}[0],
		StartTime: now.Add(-24 * time.Hour),
		EndTime:   now,
	}

	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: 1},
			Username:  "user1",
			Email:     "user1@example.com",
			Status:    1,
			Roles:     []models.Role{{Name: "admin"}},
		},
		{
			BaseModel: models.BaseModel{ID: 2},
			Username:  "user2",
			Email:     "user2@example.com",
			Status:    1,
			Roles:     []models.Role{{Name: "user"}},
		},
	}

	mockSvc.On("ExportUserList", mock.Anything, mock.AnythingOfType("*services.ExportUserListRequest")).
		Return(users, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/users/export", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["msg"])

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)
}

func TestCreateUser(t *testing.T) {
	r, mockSvc := setupTestRouter()
	r.POST("/users", CreateUser)

	// Test case 1: Successful creation
	req := services.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
		Nickname: "Test User",
		Status:   1,
	}

	createdUser := &models.User{
		BaseModel: models.BaseModel{ID: 1},
		Username:  req.Username,
		Email:     req.Email,
		Nickname:  req.Nickname,
		Status:    req.Status,
	}

	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*services.CreateUserRequest")).
		Return(createdUser, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.Username, response.Username)
	assert.Equal(t, createdUser.Email, response.Email)
}

func TestUpdateUser(t *testing.T) {
	r, mockSvc := setupTestRouter()
	r.PUT("/users/:id", UpdateUser)

	// Test case 1: Successful update
	req := services.UpdateUserRequest{
		Nickname: "Updated User",
		Email:    "updated@example.com",
		Status:   1,
	}

	updatedUser := &models.User{
		BaseModel: models.BaseModel{ID: 1},
		Username:  "testuser",
		Email:     req.Email,
		Nickname:  req.Nickname,
		Status:    req.Status,
	}

	mockSvc.On("Update", mock.Anything, uint(1), mock.AnythingOfType("*services.UpdateUserRequest")).
		Return(updatedUser, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Nickname, response.Nickname)
	assert.Equal(t, updatedUser.Email, response.Email)
}

func TestDeleteUser(t *testing.T) {
	r, mockSvc := setupTestRouter()
	r.DELETE("/users/:id", DeleteUser)

	// Test case 1: Successful deletion
	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
