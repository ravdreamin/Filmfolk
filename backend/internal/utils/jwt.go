package utils

import (
	"errors"
	"fmt"
	"time"

	"filmfolk/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// Custom JWT claims
// This is what gets encoded INTO the token
// Think of it like the "payload" of the token
type JWTClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// InitJWT initializes the JWT secret
// Call this once at app startup
func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateAccessToken creates a short-lived access token
// Access tokens are used for API requests
// They're short-lived (15 min) for security
func GenerateAccessToken(user *models.User, ttlMinutes int) (string, error) {
	if jwtSecret == nil {
		return "", errors.New("JWT secret not initialized")
	}

	// Create claims with user info and expiration
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "filmfolk",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a long-lived refresh token
// Refresh tokens are used to get new access tokens
// They're stored in the database so we can revoke them
func GenerateRefreshToken(userID uint64, ttlDays int) (string, time.Time, error) {
	if jwtSecret == nil {
		return "", time.Time{}, errors.New("JWT secret not initialized")
	}

	expiresAt := time.Now().Add(time.Duration(ttlDays) * 24 * time.Hour)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "filmfolk",
		Subject:   fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken verifies and parses a JWT token
// Returns the claims if valid, error if invalid/expired
func ValidateToken(tokenString string) (*JWTClaims, error) {
	if jwtSecret == nil {
		return nil, errors.New("JWT secret not initialized")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
// Returns the user ID if valid
func ValidateRefreshToken(tokenString string) (uint64, error) {
	if jwtSecret == nil {
		return 0, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid refresh token")
	}

	// Extract user ID from subject
	var userID uint64
	_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return 0, errors.New("invalid user ID in token")
	}

	return userID, nil
}
