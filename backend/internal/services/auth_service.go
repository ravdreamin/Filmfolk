package services

import (
	"errors"
	"fmt"
	"time"

	"filmfolk/internal/config"
	"filmfolk/internal/db"
	"filmfolk/internal/models"
	"filmfolk/internal/utils"

	"gorm.io/gorm"
)

// AuthService handles all authentication logic
type AuthService struct {
	cfg *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

// RegisterInput represents user registration data
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

// LoginInput represents user login data
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse contains tokens returned after successful auth
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds until access token expires
}

// Register creates a new user account
func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	// 1. Check if email already exists
	var existingUser models.User
	err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 2. Check if username already exists
	err = db.DB.Where("username = ?", input.Username).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("username already taken")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 3. Hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. Create user
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: &hashedPassword,
		AuthProvider: models.AuthEmail,
		Status:       models.StatusActive,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5. Generate tokens
	return s.generateAuthResponse(&user)
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	// 1. Find user by email
	var user models.User
	err := db.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 2. Check if account is active
	if user.Status != models.StatusActive {
		return nil, fmt.Errorf("account is %s", user.Status)
	}

	// 3. Verify password
	if user.PasswordHash == nil {
		return nil, errors.New("this account uses OAuth login")
	}

	if !utils.VerifyPassword(*user.PasswordHash, input.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 4. Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	db.DB.Model(&user).Update("last_login_at", now)

	// 5. Generate tokens
	return s.generateAuthResponse(&user)
}

// RefreshAccessToken generates a new access token from a refresh token
func (s *AuthService) RefreshAccessToken(refreshTokenString string) (*AuthResponse, error) {
	// 1. Validate refresh token format
	userID, err := utils.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. Check if refresh token exists in database and is valid
	var refreshToken models.RefreshToken
	err = db.DB.Where("token = ? AND user_id = ?", refreshTokenString, userID).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 3. Check if token is valid (not revoked, not expired)
	if !refreshToken.IsValid() {
		return nil, errors.New("refresh token expired or revoked")
	}

	// 4. Get user
	var user models.User
	err = db.DB.First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 5. Check if user is active
	if user.Status != models.StatusActive {
		return nil, fmt.Errorf("account is %s", user.Status)
	}

	// 6. Generate new access token (keep same refresh token)
	accessToken, err := utils.GenerateAccessToken(&user, s.cfg.Jwt.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    s.cfg.Jwt.AccessTokenTTL * 60,
	}, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(refreshTokenString string) error {
	now := time.Now()
	result := db.DB.Model(&models.RefreshToken{}).
		Where("token = ?", refreshTokenString).
		Update("revoked_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to revoke token: %w", result.Error)
	}

	return nil
}

// generateAuthResponse creates tokens and response
func (s *AuthService) generateAuthResponse(user *models.User) (*AuthResponse, error) {
	// 1. Generate access token
	accessToken, err := utils.GenerateAccessToken(user, s.cfg.Jwt.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 2. Generate refresh token
	refreshTokenString, expiresAt, err := utils.GenerateRefreshToken(user.ID, s.cfg.Jwt.RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 3. Store refresh token in database
	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}

	if err := db.DB.Create(&refreshToken).Error; err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// 4. Return response
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    s.cfg.Jwt.AccessTokenTTL * 60, // convert minutes to seconds
	}, nil
}
