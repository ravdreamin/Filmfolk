package middleware

import (
	"net/http"

	"filmfolk/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware implements rate limiting per IP address
func RateLimitMiddleware() gin.HandlerFunc {
	// Define rate: 100 requests per minute per IP
	rate := limiter.Rate{
		Period: 1 * 60 * 1000000000, // 1 minute in nanoseconds
		Limit:  100,
	}

	// Create in-memory store
	store := memory.NewStore()

	// Create limiter instance
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Get rate limit context
		context, err := instance.Get(c, ip)
		if err != nil {
			utils.GetLogger().Error().
				Err(err).
				Str("ip", ip).
				Msg("Failed to get rate limit context")
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", string(rune(context.Limit)))
		c.Header("X-RateLimit-Remaining", string(rune(context.Remaining)))
		c.Header("X-RateLimit-Reset", string(rune(context.Reset)))

		// Check if limit exceeded
		if context.Reached {
			utils.GetLogger().Warn().
				Str("ip", ip).
				Int64("limit", context.Limit).
				Msg("Rate limit exceeded")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": context.Reset,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware implements stricter rate limiting for auth endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	// Stricter rate: 10 requests per minute for auth endpoints (prevent brute force)
	rate := limiter.Rate{
		Period: 1 * 60 * 1000000000, // 1 minute
		Limit:  10,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		context, err := instance.Get(c, ip)
		if err != nil {
			utils.GetLogger().Error().
				Err(err).
				Str("ip", ip).
				Msg("Failed to get auth rate limit context")
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", string(rune(context.Limit)))
		c.Header("X-RateLimit-Remaining", string(rune(context.Remaining)))
		c.Header("X-RateLimit-Reset", string(rune(context.Reset)))

		if context.Reached {
			utils.GetLogger().Warn().
				Str("ip", ip).
				Str("endpoint", c.Request.URL.Path).
				Msg("Auth rate limit exceeded - possible brute force attempt")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many authentication attempts",
				"message": "Please wait before trying again.",
				"retry_after": context.Reset,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
