package models

import "time"

// Friendship represents the connection between two users
// This table stores bidirectional relationships
type Friendship struct {
	ID       uint64 `gorm:"primarykey" json:"id"`
	UserID   uint64 `gorm:"not null" json:"user_id"`
	FriendID uint64 `gorm:"not null" json:"friend_id"`

	// Status can be: 'pending', 'accepted', 'rejected', 'blocked'
	Status string `gorm:"size:20;not null;default:pending" json:"status"`

	// TasteSimilarityScore: Algorithm-calculated score (0-100)
	// Higher score = more similar movie taste
	// Used for friend recommendations
	TasteSimilarityScore *float64 `gorm:"type:decimal(5,2)" json:"taste_similarity_score,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Friend User `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

func (Friendship) TableName() string {
	return "friendships"
}

// FriendshipStatus constants
const (
	FriendshipPending  = "pending"
	FriendshipAccepted = "accepted"
	FriendshipRejected = "rejected"
	FriendshipBlocked  = "blocked"
)