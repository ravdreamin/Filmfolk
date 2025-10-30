package models

import "time"

type CommunityType string

const (
	CommunityPublic     CommunityType = "public"     // Anyone can join
	CommunityPrivate    CommunityType = "private"    // Invite only
	CommunityRestricted CommunityType = "restricted" // Request to join
)

// Community represents a topic-based chat room
// Example: "Horror Movie Fans", "Christopher Nolan Discussion"
type Community struct {
	ID              uint64        `gorm:"primarykey" json:"id"`
	Name            string        `gorm:"uniqueIndex;size:255;not null" json:"name"`
	Description     *string       `gorm:"type:text" json:"description,omitempty"`
	Type            CommunityType `gorm:"type:community_type;not null;default:public" json:"type"`
	CreatedByUserID uint64        `gorm:"not null" json:"created_by_user_id"`
	AvatarURL       *string       `gorm:"type:text" json:"avatar_url,omitempty"`
	MemberCount     int           `gorm:"default:0" json:"member_count"`
	IsActive        bool          `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Relationships
	CreatedBy User              `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	Members   []CommunityMember `gorm:"foreignKey:CommunityID" json:"members,omitempty"`
	Messages  []CommunityMessage `gorm:"foreignKey:CommunityID" json:"messages,omitempty"`
}

func (Community) TableName() string {
	return "communities"
}

// CommunityMember represents a user's membership in a community
type CommunityMember struct {
	ID          uint64    `gorm:"primarykey" json:"id"`
	CommunityID uint64    `gorm:"not null" json:"community_id"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	IsModerator bool      `gorm:"default:false" json:"is_moderator"` // Community-specific moderator
	JoinedAt    time.Time `gorm:"not null" json:"joined_at"`

	// Relationships
	Community Community `gorm:"foreignKey:CommunityID" json:"community,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (CommunityMember) TableName() string {
	return "community_members"
}

// CommunityMessage represents messages in a community chat
type CommunityMessage struct {
	ID              uint64    `gorm:"primarykey" json:"id"`
	CommunityID     uint64    `gorm:"not null" json:"community_id"`
	UserID          uint64    `gorm:"not null" json:"user_id"`
	MessageText     string    `gorm:"type:text;not null" json:"message_text"`
	IsRemoved       bool      `gorm:"default:false" json:"is_removed"`
	RemovedByUserID *uint64   `json:"removed_by_user_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	Community Community `gorm:"foreignKey:CommunityID" json:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RemovedBy *User     `gorm:"foreignKey:RemovedByUserID" json:"removed_by,omitempty"`
}

func (CommunityMessage) TableName() string {
	return "community_messages"
}