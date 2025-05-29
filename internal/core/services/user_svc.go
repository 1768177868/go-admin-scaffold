package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"app/internal/core/models"
	"app/internal/core/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUsernameTaken     = errors.New("username is already taken")
	ErrEmailTaken        = errors.New("email is already taken")
	ErrInvalidUserStatus = errors.New("invalid user status")
)

type UserService struct {
	userRepo *repositories.UserRepository
	logSvc   *LogService
}

func NewUserService(userRepo *repositories.UserRepository, logSvc *LogService) *UserService {
	return &UserService{
		userRepo: userRepo,
		logSvc:   logSvc,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ExportUserListRequest represents the request parameters for exporting user list
type ExportUserListRequest struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    *int      `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
	// Check if username exists
	if _, err := s.userRepo.FindByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameTaken
	}

	// Check if email exists
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailTaken
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Record operation log
	if s.logSvc != nil {
		s.logSvc.RecordOperationLog(ctx, &models.OperationLog{
			UserID:       user.ID,
			Username:     user.Username,
			Action:       "create_user",
			Module:       "user",
			BusinessID:   strconv.FormatUint(uint64(user.ID), 10),
			BusinessType: "user",
			Status:       1,
			ErrorMessage: "",
		})
	}

	return user, nil
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, id uint, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if email is taken by another user
	if req.Email != "" && req.Email != user.Email {
		if existingUser, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil && existingUser.ID != id {
			return nil, ErrEmailTaken
		}
	}

	// Update fields
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != 0 {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Record operation log
	if s.logSvc != nil {
		s.logSvc.RecordOperationLog(ctx, &models.OperationLog{
			UserID:       user.ID,
			Username:     user.Username,
			Action:       "update_user",
			Module:       "user",
			BusinessID:   strconv.FormatUint(uint64(user.ID), 10),
			BusinessType: "user",
			Status:       1,
			ErrorMessage: "",
		})
	}

	return user, nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id uint) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Record operation log
	if s.logSvc != nil {
		s.logSvc.RecordOperationLog(ctx, &models.OperationLog{
			UserID:       user.ID,
			Username:     user.Username,
			Action:       "delete_user",
			Module:       "user",
			BusinessID:   strconv.FormatUint(uint64(user.ID), 10),
			BusinessType: "user",
			Status:       1,
			ErrorMessage: "",
		})
	}

	return nil
}

// GetByID gets a user by ID
func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// List gets a paginated list of users
func (s *UserService) List(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	return s.userRepo.ListWithRoles(ctx, pagination)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, id uint, req *ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Record operation log
	if s.logSvc != nil {
		s.logSvc.RecordOperationLog(ctx, &models.OperationLog{
			UserID:       user.ID,
			Username:     user.Username,
			Action:       "change_password",
			Module:       "user",
			BusinessID:   strconv.FormatUint(uint64(user.ID), 10),
			BusinessType: "user",
			Status:       1,
			ErrorMessage: "",
		})
	}

	return nil
}

// UpdateStatus updates a user's status
func (s *UserService) UpdateStatus(ctx context.Context, id uint, status int) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if status != 0 && status != 1 {
		return ErrInvalidUserStatus
	}

	user.Status = status
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Record operation log
	if s.logSvc != nil {
		s.logSvc.RecordOperationLog(ctx, &models.OperationLog{
			UserID:       user.ID,
			Username:     user.Username,
			Action:       "update_status",
			Module:       "user",
			BusinessID:   strconv.FormatUint(uint64(user.ID), 10),
			BusinessType: "user",
			Status:       1,
			ErrorMessage: "",
		})
	}

	return nil
}

// GetUserLoginHistory gets a user's recent login history
func (s *UserService) GetUserLoginHistory(ctx context.Context, id uint, limit int) ([]models.LoginLog, error) {
	if s.logSvc == nil {
		return nil, nil
	}
	return s.logSvc.GetUserLoginHistory(ctx, id, limit)
}

// GetUserOperationHistory gets a user's recent operation history
func (s *UserService) GetUserOperationHistory(ctx context.Context, id uint, limit int) ([]models.OperationLog, error) {
	if s.logSvc == nil {
		return nil, nil
	}
	return s.logSvc.GetUserOperationHistory(ctx, id, limit)
}

// ExportUserList exports user list data based on filter criteria
func (s *UserService) ExportUserList(ctx context.Context, req *ExportUserListRequest) ([]models.User, error) {
	db := s.userRepo.GetDB().WithContext(ctx)

	// Apply filters
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if !req.StartTime.IsZero() {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		db = db.Where("created_at <= ?", req.EndTime)
	}

	var users []models.User
	err := db.Preload("Roles").Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Record operation log
	if s.logSvc != nil {
		s.logSvc.RecordOperationLog(ctx, &models.OperationLog{
			Action:       "export_users",
			Module:       "user",
			BusinessType: "user",
			Status:       1,
			ErrorMessage: "",
		})
	}

	return users, nil
}
