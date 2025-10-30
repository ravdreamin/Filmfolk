package models

import "time"

// WorldChatMessage represents the global public chat
// Unlike communities, this is ONE global room for ALL users
// Think of it like a massive Twitch chat or Reddit live thread
type WorldChatMessage struct {
	ID              uint64    `gorm:"primarykey" json:"id"`
	UserID          uint64    `gorm:"not null" json:"user_id"`
	MessageText     string    `gorm:"type:text;not null" json:"message_text"`
	IsRemoved       bool      `gorm:"default:false" json:"is_removed"`
	RemovedByUserID *uint64   `json:"removed_by_user_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`

	// Relationships
	User      User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RemovedBy *User `gorm:"foreignKey:RemovedByUserID" json:"removed_by,omitempty"`
}

func (WorldChatMessage) TableName() string {
	return "world_chat_messages"
}
