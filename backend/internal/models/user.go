package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string
type AuthProvider string
type AccountStatus string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"
)

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
	Role         UserRole      `gorm:"type:user_role;not null;default:user" json:"role"`
	AuthProvider AuthProvider  `gorm:"type:auth_provider;not null;default:email" json:"auth_provider"`
	ProviderID   *string       `gorm:"type:varchar(255)" json:"-"`
	Status       AccountStatus `gorm:"type:account_status;not null;default:active" json:"status"`
	AvatarURL    *string       `gorm:"type:text" json:"avatar_url,omitempty"`
	Bio          *string       `gorm:"type:text" json:"bio,omitempty"`

	// Gamification
	TotalReviews        int    `gorm:"default:0" json:"total_reviews"`
	TotalComments       int    `gorm:"default:0" json:"total_comments"`
	TotalLikesReceived  int    `gorm:"default:0" json:"total_likes_received"`
	EngagementScore     int    `gorm:"default:0" json:"engagement_score"`
	CurrentTitleID      *uint64 `json:"current_title_id,omitempty"`
	CurrentTitle        *UserTitle `gorm:"foreignKey:CurrentTitleID" json:"current_title,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`

	// Relationships
	Reviews          []Review          `gorm:"foreignKey:UserID" json:"-"`
	ReviewComments   []ReviewComment   `gorm:"foreignKey:UserID" json:"-"`
	UserLists        []UserList        `gorm:"foreignKey:UserID" json:"-"`
	SentMessages     []DirectMessage   `gorm:"foreignKey:SenderID" json:"-"`
	ReceivedMessages []DirectMessage   `gorm:"foreignKey:ReceiverID" json:"-"`
	Warnings         []UserWarning     `gorm:"foreignKey:UserID" json:"-"`
	Notifications    []Notification    `gorm:"foreignKey:UserID" json:"-"`
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

// Public user info for API responses
type UserPublic struct {
	ID              uint64     `json:"id"`
	Username        string     `json:"username"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	Bio             *string    `json:"bio,omitempty"`
	Role            UserRole   `json:"role"`
	CurrentTitle    *UserTitle `json:"current_title,omitempty"`
	EngagementScore int        `json:"engagement_score"`
	TotalReviews    int        `json:"total_reviews"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:              u.ID,
		Username:        u.Username,
		AvatarURL:       u.AvatarURL,
		Bio:             u.Bio,
		Role:            u.Role,
		CurrentTitle:    u.CurrentTitle,
		EngagementScore: u.EngagementScore,
		TotalReviews:    u.TotalReviews,
		CreatedAt:       u.CreatedAt,
	}
}