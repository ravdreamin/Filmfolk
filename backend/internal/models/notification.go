package models

import "time"

// Notification represents in-app notifications for users
// Examples: "Someone liked your review", "New friend request", etc.
type Notification struct {
	ID       uint64 `gorm:"primarykey" json:"id"`
	UserID   uint64 `gorm:"not null" json:"user_id"`
	Type     string `gorm:"size:50;not null" json:"type"` // e.g., 'review_like', 'friend_request', 'comment_reply'
	Title    string `gorm:"size:255;not null" json:"title"`
	Message  string `gorm:"type:text;not null" json:"message"`
	LinkURL  *string `gorm:"type:text" json:"link_url,omitempty"` // Where to navigate when clicked
	IsRead   bool    `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}

// Notification type constants
const (
	NotifTypeReviewLike     = "review_like"
	NotifTypeCommentReply   = "comment_reply"
	NotifTypeCommentLike    = "comment_like"
	NotifTypeFriendRequest  = "friend_request"
	NotifTypeFriendAccepted = "friend_accepted"
	NotifTypeWarning        = "warning"
	NotifTypeSuspension     = "suspension"
	NotifTypeNewFollower    = "new_follower"
)
