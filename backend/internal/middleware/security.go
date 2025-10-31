package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable browser XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Prevent clickjacking attacks
		c.Header("X-Frame-Options", "DENY")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy (formerly Feature Policy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (HTTP Strict Transport Security) - enforce HTTPS
		// Only set in production and when using HTTPS
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing with proper security
func CORSMiddleware(allowedOrigins []string, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// In development, allow all origins for easier testing
		if env == "development" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// In production, only allow specific origins
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			if !allowed && len(allowedOrigins) > 0 {
				// If origin not allowed, don't set CORS headers
				c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
