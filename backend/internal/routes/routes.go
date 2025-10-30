package routes

import (
	"filmfolk/internal/config"
	"filmfolk/internal/handlers"
	"filmfolk/internal/middleware"
	"filmfolk/internal/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	movieHandler := handlers.NewMovieHandler()
	reviewHandler := handlers.NewReviewHandler()

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes - no authentication required
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Public movie browsing (optional auth for personalization)
		movies := v1.Group("/movies")
		movies.Use(middleware.OptionalAuthMiddleware())
		{
			movies.GET("", movieHandler.ListMovies)           // List/search movies
			movies.GET("/:id", movieHandler.GetMovie)         // Get movie details
			movies.GET("/:id/reviews", reviewHandler.GetMovieReviews) // Get reviews for movie
		}

		// Public review viewing
		reviews := v1.Group("/reviews")
		reviews.Use(middleware.OptionalAuthMiddleware())
		{
			reviews.GET("/:id", reviewHandler.GetReview)      // Get single review with comments
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
				authMovies.POST("", movieHandler.CreateMovie)         // Submit new movie
				authMovies.PUT("/:id", movieHandler.UpdateMovie)      // Update movie (mod/admin)
			}

			// Review management
			authReviews := authenticated.Group("/reviews")
			{
				authReviews.POST("", reviewHandler.CreateReview)              // Create review
				authReviews.PUT("/:id", reviewHandler.UpdateReview)           // Update own review
				authReviews.DELETE("/:id", reviewHandler.DeleteReview)        // Delete own review
				authReviews.POST("/:id/lock", reviewHandler.LockThread)       // Lock review thread
				authReviews.POST("/:id/unlock", reviewHandler.UnlockThread)   // Unlock review thread
				authReviews.POST("/comments", reviewHandler.CreateComment)    // Add comment
				authReviews.DELETE("/comments/:id", reviewHandler.DeleteComment) // Delete comment
			}

			// TODO: User profile, lists, social features
			// users := authenticated.Group("/users")
			// lists := authenticated.Group("/lists")
			// friends := authenticated.Group("/friends")
			// messages := authenticated.Group("/messages")
			// communities := authenticated.Group("/communities")
		}

		// Moderator routes
		moderator := v1.Group("/moderator")
		moderator.Use(middleware.AuthMiddleware())
		moderator.Use(middleware.RequireRole(models.RoleModerator))
		{
			// Movie moderation
			moderator.GET("/movies/pending", movieHandler.GetPendingMovies)
			moderator.POST("/movies/:id/approve", movieHandler.ApproveMovie)
			moderator.POST("/movies/:id/reject", movieHandler.RejectMovie)

			// TODO: Review moderation
			// moderator.GET("/reviews/flagged", moderationHandler.GetFlaggedReviews)
			// moderator.POST("/reviews/:id/remove", moderationHandler.RemoveReview)
			// moderator.POST("/warnings", moderationHandler.IssueWarning)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RequireRole(models.RoleAdmin))
		{
			admin.DELETE("/movies/:id", movieHandler.DeleteMovie)

			// TODO: User management
			// admin.POST("/users/:id/ban", adminHandler.BanUser)
			// admin.POST("/users/:id/suspend", adminHandler.SuspendUser)
			// admin.GET("/moderation/logs", adminHandler.GetModerationLogs)
			// admin.GET("/stats", adminHandler.GetSystemStats)
		}
	}

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "filmfolk-api",
		})
	})
}
