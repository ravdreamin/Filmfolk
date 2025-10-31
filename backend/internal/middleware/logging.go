package middleware

import (
	"time"

	"filmfolk/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in headers
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context and response header
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware logs all HTTP requests with structured logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Get request ID
		requestID, _ := c.Get("requestID")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Build full path
		if raw != "" {
			path = path + "?" + raw
		}

		// Log level based on status code
		var event *zerolog.Event
		logger := utils.GetLogger()

		switch {
		case statusCode >= 500:
			event = logger.Error()
		case statusCode >= 400:
			event = logger.Warn()
		default:
			event = logger.Info()
		}

		event.
			Str("request_id", requestID.(string)).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency_ms", latency).
			Str("client_ip", clientIP).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size()).
			Msg("HTTP Request")

		// Log errors if any
		if len(c.Errors) > 0 {
			logger.Error().
				Str("request_id", requestID.(string)).
				Interface("errors", c.Errors.Errors()).
				Msg("Request completed with errors")
		}
	}
}
