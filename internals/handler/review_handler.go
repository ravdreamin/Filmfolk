package handler

import (
	"net/http"
	"strconv"

	"filmfolk/internals/db"
	"filmfolk/internals/models"

	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	MovieID    uint   `json:"movie_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,gte=1,lte=10"`
	ReviewText string `json:"review_text"`
}

func CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	review := models.Review{
		UserID:     userID.(uint), // Need to assert type from JWT claims
		MovieID:    req.MovieID,
		Rating:     req.Rating,
		ReviewText: req.ReviewText,
		Status:     models.ReviewStatusPublished, // Default status
		Sentiment:  "neutral",              // Default sentiment
	}

	// --- AI INTEGRATION POINT ---
	// 1. Call AI Content Moderation API with review.ReviewText
	//    - If flagged, set review.Status = models.StatusPendingModeration
	//
	// 2. If content is clean, call AI Sentiment Analysis API
	//    - Set review.Sentiment with the result ("positive", "negative")
	//
	// For now, we'll use the defaults.

	if err := db.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func GetReviewsForMovie(c *gin.Context) {
	movieID := c.Param("movieId")

	var reviews []models.Review
	result := db.DB.Where("movie_id = ? AND status = ?", movieID, models.ReviewStatusPublished).Preload("User").Find(&reviews)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func UpdateReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req struct {
		Rating     int    `json:"rating" binding:"required,gte=1,lte=10"`
		ReviewText string `json:"review_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	var review models.Review
	if err := db.DB.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this review"})
		return
	}

	review.Rating = req.Rating
	review.ReviewText = req.ReviewText

	if err := db.DB.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func DeleteReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	var review models.Review
	if err := db.DB.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	isAuthor := review.UserID == userID.(uint)
	isModerator := models.UserRole(userRole.(string)) == models.RoleModerator
	isAdmin := models.UserRole(userRole.(string)) == models.RoleAdmin

	if !isAuthor && !isModerator && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this review"})
		return
	}

	if err := db.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
