package services

import (
	"context"
	"testing"
	"time"

	"app/internal/core/models"
	"app/internal/core/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
	*repositories.UserRepository
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ListWithRoles(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

// MockLogService is a mock implementation of LogService
type MockLogService struct {
	mock.Mock
	*LogService
}

func (m *MockLogService) RecordOperationLog(ctx context.Context, log *models.OperationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockLogService) RecordLoginLog(ctx context.Context, userID uint, username, ip, userAgent string, status int, message string) error {
	args := m.Called(ctx, userID, username, ip, userAgent, status, message)
	return args.Error(0)
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func NewMockLogService() *MockLogService {
	return &MockLogService{}
}

func TestUserService_Create(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockLogSvc := NewMockLogService()
	userSvc := NewUserService(mockRepo, mockLogSvc)

	ctx := context.Background()
	req := &CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
		Nickname: "Test User",
		Status:   1,
	}

	// Test case 1: Successful user creation
	mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	mockLogSvc.On("RecordOperationLog", ctx, mock.AnythingOfType("*models.OperationLog")).Return(nil)

	user, err := userSvc.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)

	// Test case 2: Username already exists
	mockRepo.On("FindByUsername", ctx, "existinguser").Return(&models.User{}, nil)
	req.Username = "existinguser"
	_, err = userSvc.Create(ctx, req)
	assert.Equal(t, ErrUsernameTaken, err)
}

func TestUserService_ExportUserList(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogSvc := new(MockLogService)
	userSvc := NewUserService(mockRepo, mockLogSvc)

	ctx := context.Background()
	now := time.Now()

	testUsers := []models.User{
		{
			BaseModel: models.BaseModel{ID: 1},
			Username:  "user1",
			Email:     "user1@example.com",
			Status:    1,
			Roles: []models.Role{
				{Name: "admin"},
			},
		},
		{
			BaseModel: models.BaseModel{ID: 2},
			Username:  "user2",
			Email:     "user2@example.com",
			Status:    1,
			Roles: []models.Role{
				{Name: "user"},
			},
		},
	}

	req := &ExportUserListRequest{
		Username:  "user",
		Email:     "@example.com",
		Status:    &[]int{1}[0],
		StartTime: now.Add(-24 * time.Hour),
		EndTime:   now,
	}

	mockDB := &gorm.DB{}
	mockRepo.On("GetDB").Return(mockDB)

	// Mock the Find method to return test users
	mockRepo.On("Find", mock.Anything).Return(testUsers, nil)
	mockLogSvc.On("RecordOperationLog", ctx, mock.AnythingOfType("*models.OperationLog")).Return(nil)

	users, err := userSvc.ExportUserList(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogSvc := new(MockLogService)
	userSvc := NewUserService(mockRepo, mockLogSvc)

	ctx := context.Background()
	userID := uint(1)
	existingUser := &models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Nickname:  "Test User",
		Status:    1,
	}

	req := &UpdateUserRequest{
		Nickname: "Updated User",
		Email:    "updated@example.com",
		Status:   1,
	}

	// Test case 1: Successful update
	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	mockLogSvc.On("RecordOperationLog", ctx, mock.AnythingOfType("*models.OperationLog")).Return(nil)

	user, err := userSvc.Update(ctx, userID, req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Nickname, user.Nickname)
	assert.Equal(t, req.Email, user.Email)

	// Test case 2: User not found
	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)
	_, err = userSvc.Update(ctx, uint(999), req)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserService_Delete(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogSvc := new(MockLogService)
	userSvc := NewUserService(mockRepo, mockLogSvc)

	ctx := context.Background()
	userID := uint(1)
	existingUser := &models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "testuser",
	}

	// Test case 1: Successful deletion
	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Delete", ctx, userID).Return(nil)
	mockLogSvc.On("RecordOperationLog", ctx, mock.AnythingOfType("*models.OperationLog")).Return(nil)

	err := userSvc.Delete(ctx, userID)
	assert.NoError(t, err)

	// Test case 2: User not found
	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)
	err = userSvc.Delete(ctx, uint(999))
	assert.Equal(t, ErrUserNotFound, err)
}
