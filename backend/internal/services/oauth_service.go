package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"filmfolk/internal/config"
	"filmfolk/internal/db"
	"filmfolk/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// OAuthService handles OAuth authentication logic
type OAuthService struct {
	cfg          *config.Config
	googleConfig *oauth2.Config
}

// GoogleUserInfo represents user data from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(cfg *config.Config) *OAuthService {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.GoogleClientID,
		ClientSecret: cfg.OAuth.GoogleClientSecret,
		RedirectURL:  cfg.OAuth.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthService{
		cfg:          cfg,
		googleConfig: googleConfig,
	}
}

// GetGoogleAuthURL generates the Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleGoogleCallback processes the Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(code string) (*AuthResponse, error) {
	// 1. Exchange authorization code for token
	ctx := context.Background()
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// 2. Get user info from Google
	client := s.googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// 3. Validate email is verified
	if !googleUser.VerifiedEmail {
		return nil, errors.New("email not verified with Google")
	}

	// 4. Find or create user
	user, err := s.findOrCreateGoogleUser(&googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	// 5. Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	db.DB.Model(&user).Update("last_login_at", now)

	// 6. Generate JWT tokens
	return s.generateAuthResponse(user)
}

// findOrCreateGoogleUser finds existing user or creates new one
func (s *OAuthService) findOrCreateGoogleUser(googleUser *GoogleUserInfo) (*models.User, error) {
	var user models.User

	// Try to find user by Google provider ID
	err := db.DB.Where("auth_provider = ? AND provider_id = ?", models.AuthGoogle, googleUser.ID).
		First(&user).Error

	if err == nil {
		// User exists, check if active
		if user.Status != models.StatusActive {
			return nil, fmt.Errorf("account is %s", user.Status)
		}
		return &user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user with this email already exists (different auth provider)
	err = db.DB.Where("email = ?", googleUser.Email).First(&user).Error
	if err == nil {
		// Email exists but different provider
		return nil, errors.New("email already registered with different login method")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create new user
	username := s.generateUsername(googleUser)
	providerID := googleUser.ID

	user = models.User{
		Username:     username,
		Email:        googleUser.Email,
		PasswordHash: nil, // OAuth users don't have passwords
		AuthProvider: models.AuthGoogle,
		ProviderID:   &providerID,
		Status:       models.StatusActive,
		AvatarURL:    &googleUser.Picture,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// generateUsername creates a unique username from Google user info
func (s *OAuthService) generateUsername(googleUser *GoogleUserInfo) string {
	// Try name first
	baseUsername := googleUser.GivenName
	if baseUsername == "" {
		baseUsername = googleUser.Name
	}
	if baseUsername == "" {
		baseUsername = "user"
	}

	// Remove spaces and special characters
	baseUsername = sanitizeUsername(baseUsername)

	// Check if username exists
	username := baseUsername
	counter := 1

	for {
		var existing models.User
		err := db.DB.Where("username = ?", username).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Username available
			return username
		}

		// Try with number suffix
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++

		// Safety limit
		if counter > 1000 {
			// Fallback to random suffix
			username = fmt.Sprintf("%s%d", baseUsername, time.Now().UnixNano()%10000)
			break
		}
	}

	return username
}

// sanitizeUsername removes invalid characters from username
func sanitizeUsername(s string) string {
	var result []rune
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result = append(result, r)
		}
	}
	username := string(result)
	if len(username) == 0 {
		return "user"
	}
	if len(username) > 50 {
		return username[:50]
	}
	return username
}

// generateAuthResponse creates JWT tokens for OAuth user
func (s *OAuthService) generateAuthResponse(user *models.User) (*AuthResponse, error) {
	// Create service to reuse token generation logic
	authService := NewAuthService(s.cfg)
	return authService.generateAuthResponse(user)
}
