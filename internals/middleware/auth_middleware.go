package middleware

import (
	"filmfolk/internals/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == " " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Format"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected Signin method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil

		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
				return

			}

			// The vlue stored in the Gin context is not type-asserted to a concrete type (e.g., uint, string).
			// This can lead to runtime panics or unexpected behavior when handlers later try to retrieve
			// and cast it to the wrong type.  Always assert the claim to the expected type before storing.
			userID, ok := claims["user_id"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id claim type"})
				return
			}
			c.Set("userID", userID)
			c.Set("userRole", claims["user_role"])
			c.Next()

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		}

	}
}

func RoleAuthMiddleware(requireRole models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exist := c.Get("userRole")
		if !exist {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User role not found in token"})
			return
		}
		userRole := models.UserRole(role.(string))

		if userRole == models.RoleAdmin {
			c.Next()
			return
		}

		if userRole != requireRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Forbidden: require %s role", requireRole)})
			c.Next()
		}

	}
}
