package services

import (
	"context"
	"errors"
	"time"

	"app/internal/config"
	"app/internal/core/models"
	"app/internal/core/repositories"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
)

type AuthService struct {
	userRepo *repositories.UserRepository
	logSvc   *LogService
	config   *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, logSvc *LogService, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logSvc:   logSvc,
		config:   config,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		// Record failed login attempt
		if s.logSvc != nil {
			s.logSvc.RecordLoginLog(ctx, 0, req.Username, "", "", 0, "user not found")
		}
		return nil, ErrInvalidCredentials
	}

	if user.Status == 0 {
		// Record failed login attempt for inactive user
		if s.logSvc != nil {
			s.logSvc.RecordLoginLog(ctx, user.ID, user.Username, "", "", 0, "user is inactive")
		}
		return nil, ErrUserInactive
	}

	if !s.validatePassword(user.Password, req.Password) {
		// Record failed login attempt for invalid password
		if s.logSvc != nil {
			s.logSvc.RecordLoginLog(ctx, user.ID, user.Username, "", "", 0, "invalid password")
		}
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		// logger.Error("Failed to update last login time", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Record successful login
	if s.logSvc != nil {
		s.logSvc.RecordLoginLog(ctx, user.ID, user.Username, "", "", 1, "login successful")
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.config.JWT.ExpireTime * 3600, // Convert hours to seconds
	}, nil
}

func (s *AuthService) validatePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.JWT.ExpireTime)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// RefreshToken generates a new access token for the given user ID
func (s *AuthService) RefreshToken(ctx context.Context, userID uint) (string, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Status != 1 {
		return "", ErrUserInactive
	}

	return s.generateToken(user)
}
