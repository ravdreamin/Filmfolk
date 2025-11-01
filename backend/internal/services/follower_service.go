package services

import (
	"errors"
	"fmt"

	"filmfolk/internal/db"
	"filmfolk/internal/models"

	"gorm.io/gorm"
)

type FollowerService struct{}

func NewFollowerService() *FollowerService {
	return &FollowerService{}
}

// FollowUser creates a following relationship
func (s *FollowerService) FollowUser(followerID, followingID uint64) error {
	// Validation
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Check if following user exists
	var followingUser models.User
	if err := db.DB.First(&followingUser, followingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user to follow not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Check if already following
	var existing models.Follower
	err := db.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existing).Error
	if err == nil {
		return errors.New("already following this user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error: %w", err)
	}

	// Create follow relationship
	follow := models.Follower{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := db.DB.Create(&follow).Error; err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	return nil
}

// UnfollowUser removes a following relationship
func (s *FollowerService) UnfollowUser(followerID, followingID uint64) error {
	// Validation
	if followerID == followingID {
		return errors.New("invalid operation")
	}

	// Find and delete the follow relationship
	result := db.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follower{})

	if result.Error != nil {
		return fmt.Errorf("failed to unfollow user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("not following this user")
	}

	return nil
}

// IsFollowing checks if followerID is following followingID
func (s *FollowerService) IsFollowing(followerID, followingID uint64) (bool, error) {
	var count int64
	err := db.DB.Model(&models.Follower{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	return count > 0, nil
}

// GetFollowers returns list of users following the given user
func (s *FollowerService) GetFollowers(userID uint64, page, pageSize int) ([]models.UserPublic, int64, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Count total followers
	var total int64
	if err := db.DB.Model(&models.Follower{}).Where("following_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	// Get followers with user details
	var followers []models.User
	offset := (page - 1) * pageSize

	err := db.DB.Table("users").
		Select("users.*").
		Joins("INNER JOIN followers ON users.id = followers.follower_id").
		Where("followers.following_id = ?", userID).
		Order("followers.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&followers).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch followers: %w", err)
	}

	// Convert to public users
	publicUsers := make([]models.UserPublic, len(followers))
	for i, user := range followers {
		publicUsers[i] = user.ToPublic()
	}

	return publicUsers, total, nil
}

// GetFollowing returns list of users that the given user is following
func (s *FollowerService) GetFollowing(userID uint64, page, pageSize int) ([]models.UserPublic, int64, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Count total following
	var total int64
	if err := db.DB.Model(&models.Follower{}).Where("follower_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count following: %w", err)
	}

	// Get following with user details
	var following []models.User
	offset := (page - 1) * pageSize

	err := db.DB.Table("users").
		Select("users.*").
		Joins("INNER JOIN followers ON users.id = followers.following_id").
		Where("followers.follower_id = ?", userID).
		Order("followers.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&following).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch following: %w", err)
	}

	// Convert to public users
	publicUsers := make([]models.UserPublic, len(following))
	for i, user := range following {
		publicUsers[i] = user.ToPublic()
	}

	return publicUsers, total, nil
}

// GetFollowStats returns follower and following counts for a user
func (s *FollowerService) GetFollowStats(userID uint64) (followers int64, following int64, err error) {
	// Get followers count
	if err := db.DB.Model(&models.Follower{}).Where("following_id = ?", userID).Count(&followers).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	// Get following count
	if err := db.DB.Model(&models.Follower{}).Where("follower_id = ?", userID).Count(&following).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count following: %w", err)
	}

	return followers, following, nil
}
