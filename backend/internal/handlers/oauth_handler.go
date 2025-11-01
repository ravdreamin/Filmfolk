package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"filmfolk/internal/config"
	"filmfolk/internal/services"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth HTTP requests
type OAuthHandler struct {
	oauthService *services.OAuthService
	frontendURL  string
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(cfg *config.Config) *OAuthHandler {
	// Determine frontend URL from allowed origins
	frontendURL := "http://localhost:3000"
	if cfg.App.Env == "production" && len(cfg.App.AllowedOrigins) > 0 {
		frontendURL = cfg.App.AllowedOrigins[0]
	}

	return &OAuthHandler{
		oauthService: services.NewOAuthService(cfg),
		frontendURL:  frontendURL,
	}
}

// GoogleLogin initiates Google OAuth flow
// @Summary Start Google OAuth login
// @Description Redirects to Google for authentication
// @Tags auth
// @Success 302 {string} string "Redirect to Google"
// @Router /auth/google [get]
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state in session cookie (httpOnly, secure in production)
	c.SetCookie(
		"oauth_state",
		state,
		600, // 10 minutes
		"/",
		"",
		c.Request.URL.Scheme == "https",
		true, // httpOnly
	)

	// Redirect to Google OAuth consent screen
	authURL := h.oauthService.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles Google OAuth callback
// @Summary Handle Google OAuth callback
// @Description Processes Google OAuth response and creates/logs in user
// @Tags auth
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State for CSRF protection"
// @Success 302 {string} string "Redirect to frontend with token"
// @Failure 400 {object} gin.H
// @Router /auth/google/callback [get]
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	// 1. Verify state for CSRF protection
	stateFromQuery := c.Query("state")
	stateFromCookie, err := c.Cookie("oauth_state")
	if err != nil || stateFromQuery == "" || stateFromQuery != stateFromCookie {
		h.redirectToFrontendWithError(c, "invalid_state", "Invalid OAuth state")
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// 2. Check for OAuth errors
	if errorParam := c.Query("error"); errorParam != "" {
		errorDesc := c.Query("error_description")
		if errorDesc == "" {
			errorDesc = "OAuth authentication failed"
		}
		h.redirectToFrontendWithError(c, errorParam, errorDesc)
		return
	}

	// 3. Get authorization code
	code := c.Query("code")
	if code == "" {
		h.redirectToFrontendWithError(c, "no_code", "No authorization code received")
		return
	}

	// 4. Exchange code for user info and create/login user
	authResponse, err := h.oauthService.HandleGoogleCallback(code)
	if err != nil {
		h.redirectToFrontendWithError(c, "auth_failed", err.Error())
		return
	}

	// 5. Redirect to frontend with tokens
	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s",
		h.frontendURL,
		authResponse.AccessToken,
		authResponse.RefreshToken,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// redirectToFrontendWithError redirects to frontend error page
func (h *OAuthHandler) redirectToFrontendWithError(c *gin.Context, errorCode, errorMessage string) {
	redirectURL := fmt.Sprintf("%s/auth/error?code=%s&message=%s",
		h.frontendURL,
		errorCode,
		errorMessage,
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// generateRandomState generates a random state string for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
