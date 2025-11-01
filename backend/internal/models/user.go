package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthProvider string
type AccountStatus string

const (
	AuthEmail     AuthProvider = "email"
	AuthGoogle    AuthProvider = "google"
	AuthFacebook  AuthProvider = "facebook"
	AuthInstagram AuthProvider = "instagram"
	AuthTwitter   AuthProvider = "twitter"
	AuthGuest     AuthProvider = "guest"
)

const (
	StatusActive    AccountStatus = "active"
	StatusSuspended AccountStatus = "suspended"
	StatusBanned    AccountStatus = "banned"
)

type User struct {
	ID           uint64        `gorm:"primarykey" json:"id"`
	Username     string        `gorm:"uniqueIndex;not null" json:"username"`
	Email        string        `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash *string       `gorm:"type:text" json:"-"`
	AuthProvider AuthProvider  `gorm:"type:auth_provider;not null;default:email" json:"auth_provider"`
	ProviderID   *string       `gorm:"type:varchar(255)" json:"-"`
	Status       AccountStatus `gorm:"type:account_status;not null;default:active" json:"status"`
	AvatarURL    *string       `gorm:"type:text" json:"avatar_url,omitempty"`
	Bio          *string       `gorm:"type:text" json:"bio,omitempty"`

	// Follower counts
	FollowersCount int `gorm:"default:0" json:"followers_count"`
	FollowingCount int `gorm:"default:0" json:"following_count"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to validate data
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Email auth requires password
	if u.AuthProvider == AuthEmail && u.PasswordHash == nil {
		return gorm.ErrInvalidData
	}
	return nil
}

// UserPublic represents public user info for API responses
type UserPublic struct {
	ID             uint64    `json:"id"`
	Username       string    `json:"username"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	Bio            *string   `json:"bio,omitempty"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	CreatedAt      time.Time `json:"created_at"`
}

// ToPublic converts User to UserPublic
func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:             u.ID,
		Username:       u.Username,
		AvatarURL:      u.AvatarURL,
		Bio:            u.Bio,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		CreatedAt:      u.CreatedAt,
	}
}
