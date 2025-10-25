package routes

import (
	"filmfolk/internals/handler"
	"filmfolk/internals/middleware"

	"github.com/gin-gonic/gin"
)

// SetupReviewRoutes configures the review-related routes.
func SetupReviewRoutes(router *gin.RouterGroup) {
	// Create a new group for /reviews
	reviewGroup := router.Group("/reviews")
	{
		// Public route: Get all reviews for a specific movie
		// GET /api/v1/reviews/movie/:movieId
		reviewGroup.GET("/movie/:movieId", handler.GetReviewsForMovie)

		// --- Protected Routes ---
		// All routes below this point will require a valid JWT
		protected := reviewGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// POST /api/v1/reviews
			protected.POST("/", handler.CreateReview)

			// PUT /api/v1/reviews/:id
			protected.PUT("/:id", handler.UpdateReview)

			// DELETE /api/v1/reviews/:id
			protected.DELETE("/:id", handler.DeleteReview)
		}
	}
}
