package models

import "time"

type ModerationAction string

const (
	ActionReviewFlagged  ModerationAction = "review_flagged"
	ActionReviewRemoved  ModerationAction = "review_removed"
	ActionUserWarned     ModerationAction = "user_warned"
	ActionUserSuspended  ModerationAction = "user_suspended"
	ActionUserBanned     ModerationAction = "user_banned"
)

// ModerationLog tracks ALL moderation actions
// This is your audit trail - who did what, when, and why
type ModerationLog struct {
	ID           uint64           `gorm:"primarykey" json:"id"`
	ModeratorID  uint64           `gorm:"not null" json:"moderator_id"`
	TargetUserID *uint64          `json:"target_user_id,omitempty"`
	Action       ModerationAction `gorm:"type:moderation_action;not null" json:"action"`
	Reason       string           `gorm:"type:text;not null" json:"reason"`

	// References to moderated content
	ReviewID  *uint64 `json:"review_id,omitempty"`
	CommentID *uint64 `json:"comment_id,omitempty"`

	// Suspension/Ban details
	DurationDays *int       `json:"duration_days,omitempty"` // NULL = permanent
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Moderator  User           `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
	TargetUser *User          `gorm:"foreignKey:TargetUserID" json:"target_user,omitempty"`
	Review     *Review        `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
	Comment    *ReviewComment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
}

func (ModerationLog) TableName() string {
	return "moderation_logs"
}

// UserWarning tracks warnings issued to users
// Three strikes = moderator escalates to admin
type UserWarning struct {
	ID                 uint64  `gorm:"primarykey" json:"id"`
	UserID             uint64  `gorm:"not null" json:"user_id"`
	IssuedByModeratorID uint64  `gorm:"not null" json:"issued_by_moderator_id"`
	Reason             string  `gorm:"type:text;not null" json:"reason"`
	ModerationLogID    *uint64 `json:"moderation_log_id,omitempty"`
	IsActive           bool    `gorm:"default:true" json:"is_active"` // Can be revoked
	CreatedAt          time.Time `json:"created_at"`

	// Relationships
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IssuedBy    User           `gorm:"foreignKey:IssuedByModeratorID" json:"issued_by,omitempty"`
	ModerationLog *ModerationLog `gorm:"foreignKey:ModerationLogID" json:"moderation_log,omitempty"`
}

func (UserWarning) TableName() string {
	return "user_warnings"
}
