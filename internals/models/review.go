package models

import (
	"time"
)

type ReviewStatus string

const (
    StatusPendingModeration ReviewStatus = "pending_moderation"
    StatusPublished         ReviewStatus = "published"
    StatusRejected          ReviewStatus = "rejected"
)

// Review represents the reviews table and its relationships.
type Review struct {
    ID          uint64       `gorm:"primaryKey"`
    UserID      uint64       `gorm:"uniqueIndex:idx_user_movie;not null"` // Foreign key for User
    MovieID     uint64       `gorm:"uniqueIndex:idx_user_movie;not null"` // Foreign key for Movie
    Rating      int          `gorm:"not null"`
    ReviewText  string       `gorm:"not null"`
    Sentiment   string
    Status      ReviewStatus `gorm:"type:review_status;not null;default:published"`
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Define the relationships (associations) to other models.
    User  User  `gorm:"foreignKey:UserID"`
    Movie Movie `gorm:"foreignKey:MovieID"`
}
