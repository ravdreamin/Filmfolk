package middleware

import (
	"net/http"
	"strings"

	"filmfolk/internal/models"
	"filmfolk/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and adds user info to context
// This is the gatekeeper for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract token from Authorization header
		// Format: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 2. Split "Bearer" and token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Validate and parse token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Add user info to context for handlers to use
		// Handlers can now access: c.Get("userID"), c.Get("userRole"), etc.
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("userRole", claims.Role)

		// 5. Continue to next handler
		c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate but doesn't block if token is missing
// Useful for routes that change behavior based on whether user is logged in
// Example: Guest users can view reviews, but can't like them
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token, but that's okay - continue as guest
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, but don't block - continue as guest
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, but don't block - continue as guest
			c.Next()
			return
		}

		// Valid token - add user info to context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
// Must be used AFTER AuthMiddleware
// Example: RequireRole(models.RoleModerator) - only moderators/admins can access
func RequireRole(minRole models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRoleValue, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userRole, ok := userRoleValue.(models.UserRole)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		// Check role hierarchy: user < moderator < admin
		if !hasRequiredRole(userRole, minRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole checks role hierarchy
func hasRequiredRole(userRole, requiredRole models.UserRole) bool {
	roleHierarchy := map[models.UserRole]int{
		models.RoleUser:      1,
		models.RoleModerator: 2,
		models.RoleAdmin:     3,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// GetUserID is a helper to extract user ID from context
// Returns 0 if not authenticated
func GetUserID(c *gin.Context) uint64 {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint64); ok {
			return id
		}
	}
	return 0
}

// GetUserRole is a helper to extract user role from context
func GetUserRole(c *gin.Context) models.UserRole {
	if role, exists := c.Get("userRole"); exists {
		if r, ok := role.(models.UserRole); ok {
			return r
		}
	}
	return models.RoleUser
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("userID")
	return exists
}
