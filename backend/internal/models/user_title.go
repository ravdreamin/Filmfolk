package models

import "time"

type UserTitle struct {
	ID                      uint64    `gorm:"primarykey" json:"id"`
	Name                    string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description             *string   `gorm:"type:text" json:"description,omitempty"`
	RequiredReviews         int       `gorm:"default:0" json:"required_reviews"`
	RequiredComments        int       `gorm:"default:0" json:"required_comments"`
	RequiredEngagementScore int       `gorm:"default:0" json:"required_engagement_score"`
	IconURL                 *string   `gorm:"type:text" json:"icon_url,omitempty"`
	SortOrder               int       `gorm:"default:0" json:"sort_order"`
	CreatedAt               time.Time `json:"created_at"`
}

func (UserTitle) TableName() string {
	return "user_titles"
}
