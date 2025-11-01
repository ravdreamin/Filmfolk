package services

import (
	"errors"
	"fmt"

	"filmfolk/internal/db"
	"filmfolk/internal/models"

	"gorm.io/gorm"
)

// ReviewService handles review-related business logic
type ReviewService struct{}

// NewReviewService creates a new review service
func NewReviewService() *ReviewService {
	return &ReviewService{}
}

// CreateReviewInput represents data for creating a review
type CreateReviewInput struct {
	MovieID    uint64 `json:"movie_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=10"`
	ReviewText string `json:"review_text" binding:"required,min=10"`
}

// UpdateReviewInput represents data for updating a review
type UpdateReviewInput struct {
	Rating     *int    `json:"rating,omitempty"`
	ReviewText *string `json:"review_text,omitempty"`
}

// CreateCommentInput represents data for creating a comment
type CreateCommentInput struct {
	ReviewID        uint64  `json:"review_id" binding:"required"`
	ParentCommentID *uint64 `json:"parent_comment_id,omitempty"`
	CommentText     string  `json:"comment_text" binding:"required,min=1"`
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(input CreateReviewInput, userID uint64) (*models.Review, error) {
	// Check if user already reviewed this movie
	var existing models.Review
	err := db.DB.Where("user_id = ? AND movie_id = ?", userID, input.MovieID).First(&existing).Error
	if err == nil {
		return nil, errors.New("you have already reviewed this movie")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify movie exists
	var movie models.Movie
	if err := db.DB.First(&movie, input.MovieID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create review
	review := models.Review{
		UserID:     userID,
		MovieID:    input.MovieID,
		Rating:     input.Rating,
		ReviewText: input.ReviewText,
	}

	if err := db.DB.Create(&review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Update movie stats
	movieService := NewMovieService()
	movieService.RecalculateMovieStats(input.MovieID)

	// Load user relation
	db.DB.Preload("User").First(&review, review.ID)

	return &review, nil
}

// GetReview retrieves a review by ID
func (s *ReviewService) GetReview(reviewID uint64) (*models.Review, error) {
	var review models.Review
	err := db.DB.Preload("User").
		Preload("Movie").
		Preload("Comments.User").
		Preload("Comments.Replies.User").
		First(&review, reviewID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &review, nil
}

// GetReviewsForMovie retrieves all reviews for a movie
func (s *ReviewService) GetReviewsForMovie(movieID uint64, page, pageSize int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	var total int64
	query := db.DB.Model(&models.Review{}).
		Where("movie_id = ? AND status = ?", movieID, models.ReviewStatusPublished)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []models.Review
	offset := (page - 1) * pageSize
	err := db.DB.Where("movie_id = ? AND status = ?", movieID, models.ReviewStatusPublished).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return reviews, total, nil
}

// GetUserReviews retrieves all reviews by a user
func (s *ReviewService) GetUserReviews(userID uint64, page, pageSize int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	var total int64
	query := db.DB.Model(&models.Review{}).
		Where("user_id = ? AND status = ?", userID, models.ReviewStatusPublished)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []models.Review
	offset := (page - 1) * pageSize
	err := db.DB.Where("user_id = ? AND status = ?", userID, models.ReviewStatusPublished).
		Preload("Movie").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return reviews, total, nil
}

// UpdateReview updates a review
func (s *ReviewService) UpdateReview(reviewID, userID uint64, input UpdateReviewInput) (*models.Review, error) {
	var review models.Review
	if err := db.DB.First(&review, reviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check ownership
	if review.UserID != userID {
		return nil, errors.New("you can only edit your own reviews")
	}

	// Update fields
	updates := make(map[string]interface{})
	if input.Rating != nil {
		if *input.Rating < 1 || *input.Rating > 10 {
			return nil, errors.New("rating must be between 1 and 10")
		}
		updates["rating"] = *input.Rating
	}
	if input.ReviewText != nil {
		updates["review_text"] = *input.ReviewText
		// TODO: Re-run AI moderation
	}

	if err := db.DB.Model(&review).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Recalculate movie stats if rating changed
	if input.Rating != nil {
		movieService := NewMovieService()
		movieService.RecalculateMovieStats(review.MovieID)
	}

	// Reload review
	db.DB.Preload("User").Preload("Movie").First(&review, reviewID)

	return &review, nil
}

// DeleteReview deletes a review (only review owner)
func (s *ReviewService) DeleteReview(reviewID, userID uint64) error {
	var review models.Review
	if err := db.DB.First(&review, reviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Check permission (only owner can delete)
	if review.UserID != userID {
		return errors.New("you don't have permission to delete this review")
	}

	movieID := review.MovieID
	if err := db.DB.Delete(&review).Error; err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	// Update movie stats
	movieService := NewMovieService()
	movieService.RecalculateMovieStats(movieID)

	return nil
}

// LockReviewThread locks a review thread (only review author can lock)
func (s *ReviewService) LockReviewThread(reviewID, userID uint64) error {
	var review models.Review
	if err := db.DB.First(&review, reviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if review.UserID != userID {
		return errors.New("only review author can lock the thread")
	}

	if review.IsThreadLocked {
		return errors.New("thread is already locked")
	}

	return db.DB.Model(&review).Update("is_thread_locked", true).Error
}

// UnlockReviewThread unlocks a review thread
func (s *ReviewService) UnlockReviewThread(reviewID, userID uint64) error {
	var review models.Review
	if err := db.DB.First(&review, reviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if review.UserID != userID {
		return errors.New("only review author can unlock the thread")
	}

	if !review.IsThreadLocked {
		return errors.New("thread is not locked")
	}

	return db.DB.Model(&review).Update("is_thread_locked", false).Error
}

// CreateComment creates a comment on a review
func (s *ReviewService) CreateComment(input CreateCommentInput, userID uint64) (*models.ReviewComment, error) {
	// Check if review exists
	var review models.Review
	if err := db.DB.First(&review, input.ReviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if thread is locked
	if review.IsThreadLocked {
		return nil, errors.New("this review thread is locked")
	}

	// If replying to a comment, verify it exists and belongs to the same review
	if input.ParentCommentID != nil {
		var parentComment models.ReviewComment
		if err := db.DB.First(&parentComment, *input.ParentCommentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent comment not found")
			}
			return nil, fmt.Errorf("database error: %w", err)
		}

		if parentComment.ReviewID != input.ReviewID {
			return nil, errors.New("parent comment belongs to a different review")
		}
	}

	// Create comment
	comment := models.ReviewComment{
		ReviewID:        input.ReviewID,
		UserID:          userID,
		ParentCommentID: input.ParentCommentID,
		CommentText:     input.CommentText,
	}

	// TODO: AI content moderation

	if err := db.DB.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Update review comment count
	db.DB.Model(&review).UpdateColumn("comments_count", gorm.Expr("comments_count + 1"))

	// Load user relation
	db.DB.Preload("User").First(&comment, comment.ID)

	return &comment, nil
}

// DeleteComment deletes a comment (only comment owner)
func (s *ReviewService) DeleteComment(commentID, userID uint64) error {
	var comment models.ReviewComment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Check permission (only owner can delete)
	if comment.UserID != userID {
		return errors.New("you don't have permission to delete this comment")
	}

	reviewID := comment.ReviewID
	if err := db.DB.Delete(&comment).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Update review comment count
	db.DB.Model(&models.Review{}).Where("id = ?", reviewID).
		UpdateColumn("comments_count", gorm.Expr("comments_count - 1"))

	return nil
}
