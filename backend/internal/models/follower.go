package models

import "time"

// Follower represents the following relationship between users
// follower_id follows following_id
type Follower struct {
	ID          uint64    `gorm:"primarykey" json:"id"`
	FollowerID  uint64    `gorm:"not null;uniqueIndex:idx_follower_following" json:"follower_id"`  // User who is following
	FollowingID uint64    `gorm:"not null;uniqueIndex:idx_follower_following" json:"following_id"` // User being followed
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

func (Follower) TableName() string {
	return "followers"
}
