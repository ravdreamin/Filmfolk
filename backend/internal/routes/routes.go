package routes

import (
	"filmfolk/internal/config"
	"filmfolk/internal/handlers"
	"filmfolk/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	oauthHandler := handlers.NewOAuthHandler(cfg)
	movieHandler := handlers.NewMovieHandler()
	reviewHandler := handlers.NewReviewHandler()
	followerHandler := handlers.NewFollowerHandler()
	healthHandler := handlers.NewHealthHandler()

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes - no authentication required
		auth := v1.Group("/auth")
		auth.Use(middleware.AuthRateLimitMiddleware()) // Stricter rate limiting for auth
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)

			// Google OAuth routes
			auth.GET("/google", oauthHandler.GoogleLogin)
			auth.GET("/google/callback", oauthHandler.GoogleCallback)
		}

		// Public movie browsing (optional auth for personalization)
		movies := v1.Group("/movies")
		movies.Use(middleware.OptionalAuthMiddleware())
		{
			movies.GET("", movieHandler.ListMovies)                       // List/search movies
			movies.GET("/:id", movieHandler.GetMovie)                     // Get movie details
			movies.GET("/:id/reviews", reviewHandler.GetMovieReviews)     // Get reviews for movie
		}

		// Public review viewing
		reviews := v1.Group("/reviews")
		reviews.Use(middleware.OptionalAuthMiddleware())
		{
			reviews.GET("/:id", reviewHandler.GetReview) // Get single review with comments
		}

		// Public user profile routes
		users := v1.Group("/users")
		{
			users.GET("/:id/followers", followerHandler.GetFollowers)  // Get user's followers
			users.GET("/:id/following", followerHandler.GetFollowing)  // Get users that user follows
		}

		// Protected routes - authentication required
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// Current user info
			authenticated.GET("/auth/me", authHandler.GetCurrentUser)

			// Authenticated movie operations
			authMovies := authenticated.Group("/movies")
			{
				authMovies.PUT("/:id", movieHandler.UpdateMovie) // Update movie
			}

			// Review management
			authReviews := authenticated.Group("/reviews")
			{
				authReviews.POST("", reviewHandler.CreateReview)                     // Create review
				authReviews.PUT("/:id", reviewHandler.UpdateReview)                  // Update own review
				authReviews.DELETE("/:id", reviewHandler.DeleteReview)               // Delete own review
				authReviews.POST("/:id/lock", reviewHandler.LockThread)              // Lock review thread
				authReviews.POST("/:id/unlock", reviewHandler.UnlockThread)          // Unlock review thread
				authReviews.POST("/comments", reviewHandler.CreateComment)           // Add comment
				authReviews.DELETE("/comments/:id", reviewHandler.DeleteComment)     // Delete comment
			}

			// Follower management
			authUsers := authenticated.Group("/users")
			{
				authUsers.POST("/:id/follow", followerHandler.FollowUser)            // Follow a user
				authUsers.DELETE("/:id/follow", followerHandler.UnfollowUser)        // Unfollow a user
				authUsers.GET("/:id/follow/status", followerHandler.CheckFollowStatus) // Check if following
			}
		}
	}

	// Health check endpoints (no auth required, no rate limiting)
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/detailed", healthHandler.DetailedHealthCheck)
	router.GET("/health/ready", healthHandler.ReadinessCheck)
	router.GET("/health/live", healthHandler.LivenessCheck)
}
