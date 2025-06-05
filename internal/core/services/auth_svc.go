package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"app/internal/config"
	"app/internal/core/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
)

type AuthService struct {
	userRepo UserRepository
	logSvc   *LogService
	config   *config.Config
}

func NewAuthService(userRepo UserRepository, logSvc *LogService, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logSvc:   logSvc,
		config:   config,
	}
}

// IsSuperAdmin checks if a user ID is in the super admin list
func (s *AuthService) IsSuperAdmin(userID uint) bool {
	if s.config == nil {
		log.Printf("[ERROR] Config is nil when checking super admin for user %d", userID)
		return false
	}
	log.Printf("[DEBUG] Checking if user %d is super admin", userID)
	superAdminIDs := s.config.ParseSuperAdminIDs()
	log.Printf("[DEBUG] Super admin IDs from config: %v", superAdminIDs)
	for _, id := range superAdminIDs {
		if id == userID {
			log.Printf("[DEBUG] User %d is super admin", userID)
			return true
		}
	}
	log.Printf("[DEBUG] User %d is not super admin", userID)
	return false
}

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ValidateToken validates a JWT token and returns its claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// GetUserFromClaims retrieves user information from JWT claims
func (s *AuthService) GetUserFromClaims(ctx context.Context, claims jwt.MapClaims) (*models.User, error) {
	// Get user_id from claims with type checking
	userIDValue, exists := claims["user_id"]
	if !exists {
		log.Printf("[ERROR] user_id not found in claims")
		return nil, errors.New("user_id not found in claims")
	}

	var userID uint
	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case float32:
		userID = uint(v)
	case int:
		userID = uint(v)
	case int64:
		userID = uint(v)
	case uint:
		userID = v
	case uint64:
		userID = uint(v)
	default:
		log.Printf("[ERROR] invalid user_id type in claims: %T", userIDValue)
		return nil, errors.New("invalid user_id type in claims")
	}

	log.Printf("[DEBUG] Getting user from claims, user_id: %d", userID)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("[ERROR] Failed to find user by ID %d: %v", userID, err)
		return nil, err
	}

	// Set IsSuperAdmin field
	user.IsSuperAdmin = s.IsSuperAdmin(user.ID)
	log.Printf("[DEBUG] User %d IsSuperAdmin: %v", user.ID, user.IsSuperAdmin)

	return user, nil
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

	// Set IsSuperAdmin field
	user.IsSuperAdmin = s.IsSuperAdmin(user.ID)

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		log.Printf("[WARN] Failed to update last login time: %v", err)
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
		ExpiresIn:   s.config.JWT.ExpireTime, // ExpireTime is already in seconds
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
		"exp":      time.Now().Add(time.Second * time.Duration(s.config.JWT.ExpireTime)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		log.Printf("[ERROR] Failed to sign JWT token: %v", err)
		return "", err
	}

	log.Printf("[DEBUG] Generated JWT token for user %d: %s", user.ID, tokenString)
	log.Printf("[DEBUG] JWT token length: %d", len(tokenString))
	log.Printf("[DEBUG] JWT token segments: %d", len(strings.Split(tokenString, ".")))

	return tokenString, nil
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

	// Set IsSuperAdmin field
	user.IsSuperAdmin = s.IsSuperAdmin(user.ID)

	return s.generateToken(user)
}

// GetConfig returns the JWT configuration
func (s *AuthService) GetConfig() *config.Config {
	return s.config
}

// Logout handles user logout and logs the action
func (s *AuthService) Logout(ctx context.Context, userID uint) error {
	// Get user information for logging
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Record logout in login logs
	if s.logSvc != nil {
		return s.logSvc.RecordLoginLog(ctx, user.ID, user.Username, "", "", 1, "logout successful")
	}

	return nil
}
