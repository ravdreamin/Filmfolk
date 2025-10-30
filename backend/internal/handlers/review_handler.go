package handlers

import (
	"net/http"
	"strconv"

	"filmfolk/internal/middleware"
	"filmfolk/internal/services"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService: services.NewReviewService(),
	}
}

// CreateReview handles POST /api/v1/reviews
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var input services.CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	review, err := h.reviewService.CreateReview(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetReview handles GET /api/v1/reviews/:id
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.reviewService.GetReview(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// GetMovieReviews handles GET /api/v1/movies/:id/reviews
func (h *ReviewHandler) GetMovieReviews(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.reviewService.GetReviewsForMovie(movieID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":   reviews,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateReview handles PUT /api/v1/reviews/:id
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var input services.UpdateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	review, err := h.reviewService.UpdateReview(id, userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// DeleteReview handles DELETE /api/v1/reviews/:id
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	if err := h.reviewService.DeleteReview(id, userID, userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// LockThread handles POST /api/v1/reviews/:id/lock
func (h *ReviewHandler) LockThread(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.reviewService.LockReviewThread(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thread locked successfully"})
}

// UnlockThread handles POST /api/v1/reviews/:id/unlock
func (h *ReviewHandler) UnlockThread(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.reviewService.UnlockReviewThread(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thread unlocked successfully"})
}

// CreateComment handles POST /api/v1/reviews/comments
func (h *ReviewHandler) CreateComment(c *gin.Context) {
	var input services.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	comment, err := h.reviewService.CreateComment(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// DeleteComment handles DELETE /api/v1/reviews/comments/:id
func (h *ReviewHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	if err := h.reviewService.DeleteComment(id, userID, userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
