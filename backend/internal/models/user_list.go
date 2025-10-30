package models

import "time"

// ListType defines the category of user lists
// Example: 'watched', 'plan_to_watch', 'custom'
type ListType string

const (
	ListTypeWatched      ListType = "watched"
	ListTypeDropped      ListType = "dropped"
	ListTypePlanToWatch  ListType = "plan_to_watch"
	ListTypeCustom       ListType = "custom"
)

// UserList represents a collection of movies curated by a user
// Think of it like a Spotify playlist, but for movies
type UserList struct {
	ID          uint64   `gorm:"primarykey" json:"id"`
	UserID      uint64   `gorm:"not null" json:"user_id"`
	Name        string   `gorm:"size:255;not null" json:"name"`
	ListType    ListType `gorm:"type:list_type;not null" json:"list_type"`
	Description *string  `gorm:"type:text" json:"description,omitempty"`
	IsPublic    bool     `gorm:"default:true" json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	User  User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items []UserListItem `gorm:"foreignKey:ListID" json:"items,omitempty"`
}

func (UserList) TableName() string {
	return "user_lists"
}

// UserListItem represents a movie in a user's list
// The junction table connecting lists and movies
type UserListItem struct {
	ID      uint64    `gorm:"primarykey" json:"id"`
	ListID  uint64    `gorm:"not null" json:"list_id"`
	MovieID uint64    `gorm:"not null" json:"movie_id"`
	Notes   *string   `gorm:"type:text" json:"notes,omitempty"`
	AddedAt time.Time `gorm:"not null" json:"added_at"`

	// Relationships
	List  UserList `gorm:"foreignKey:ListID" json:"-"`
	Movie Movie    `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
}

func (UserListItem) TableName() string {
	return "user_list_items"
}