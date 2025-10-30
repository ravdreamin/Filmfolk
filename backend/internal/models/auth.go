package models

import "time"

// RefreshToken stores JWT refresh tokens
// Why separate table? To allow token revocation and rotation
// Access tokens are stateless (stored in memory), refresh tokens need persistence
type RefreshToken struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	Token     string    `gorm:"type:text;uniqueIndex;not null" json:"-"` // Never expose in API
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"` // NULL = still valid

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid checks if the token is still usable
func (rt *RefreshToken) IsValid() bool {
	// Token is valid if:
	// 1. Not revoked (RevokedAt is NULL)
	// 2. Not expired (ExpiresAt is in the future)
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}
