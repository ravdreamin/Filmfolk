package models

import "time"

type MessageStatus string

const (
	MessageSent      MessageStatus = "sent"
	MessageDelivered MessageStatus = "delivered"
	MessageRead      MessageStatus = "read"
)

// DirectMessage represents one-on-one messaging between users
// Similar to WhatsApp or Instagram DMs
type DirectMessage struct {
	ID          uint64        `gorm:"primarykey" json:"id"`
	SenderID    uint64        `gorm:"not null" json:"sender_id"`
	ReceiverID  uint64        `gorm:"not null" json:"receiver_id"`
	MessageText string        `gorm:"type:text;not null" json:"message_text"`
	Status      MessageStatus `gorm:"type:message_status;not null;default:sent" json:"status"`

	// Soft delete flags - message exists but hidden from one/both users
	IsDeletedBySender   bool `gorm:"default:false" json:"is_deleted_by_sender"`
	IsDeletedByReceiver bool `gorm:"default:false" json:"is_deleted_by_receiver"`

	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`

	// Relationships
	Sender   User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
}

func (DirectMessage) TableName() string {
	return "direct_messages"
}
