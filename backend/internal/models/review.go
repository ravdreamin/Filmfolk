package models

import (
	"time"
)

type ReviewStatus string

const (
	ReviewStatusPending   ReviewStatus = "pending_moderation"
	ReviewStatusPublished ReviewStatus = "published"
	ReviewStatusRejected  ReviewStatus = "rejected"
)

type Review struct {
	ID         uint64 `gorm:"primarykey" json:"id"`
	UserID     uint64 `gorm:"not null" json:"user_id"`
	MovieID    uint64 `gorm:"not null" json:"movie_id"`
	Rating     int    `gorm:"not null;check:rating >= 1 AND rating <= 10" json:"rating"`
	ReviewText string `gorm:"type:text;not null" json:"review_text"`

	// AI Analysis
	Sentiment    *string `gorm:"size:50" json:"sentiment,omitempty"`
	AIFlagged    bool    `gorm:"default:false" json:"ai_flagged"`
	AIFlagReason *string `gorm:"type:text" json:"ai_flag_reason,omitempty"`

	Status ReviewStatus `gorm:"type:review_status;not null;default:published" json:"status"`

	// Thread control
	IsThreadLocked bool `gorm:"default:false" json:"is_thread_locked"`

	// Engagement
	LikesCount    int `gorm:"default:0" json:"likes_count"`
	CommentsCount int `gorm:"default:0" json:"comments_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User     User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie    Movie           `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Comments []ReviewComment `gorm:"foreignKey:ReviewID" json:"comments,omitempty"`
	Likes    []ReviewLike    `gorm:"foreignKey:ReviewID" json:"-"`
}

func (Review) TableName() string {
	return "reviews"
}

type ReviewComment struct {
	ID              uint64  `gorm:"primarykey" json:"id"`
	ReviewID        uint64  `gorm:"not null" json:"review_id"`
	UserID          uint64  `gorm:"not null" json:"user_id"`
	ParentCommentID *uint64 `json:"parent_comment_id,omitempty"`
	CommentText     string  `gorm:"type:text;not null" json:"comment_text"`

	// Moderation
	AIFlagged       bool    `gorm:"default:false" json:"ai_flagged"`
	IsRemoved       bool    `gorm:"default:false" json:"is_removed"`
	RemovedByUserID *uint64 `json:"removed_by_user_id,omitempty"`

	LikesCount int `gorm:"default:0" json:"likes_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Review        Review           `gorm:"foreignKey:ReviewID" json:"-"`
	User          User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentComment *ReviewComment   `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	Replies       []ReviewComment  `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
	RemovedBy     *User            `gorm:"foreignKey:RemovedByUserID" json:"removed_by,omitempty"`
	Likes         []CommentLike    `gorm:"foreignKey:CommentID" json:"-"`
}

func (ReviewComment) TableName() string {
	return "review_comments"
}

type ReviewLike struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	ReviewID  uint64    `gorm:"not null" json:"review_id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Review Review `gorm:"foreignKey:ReviewID" json:"-"`
	User   User   `gorm:"foreignKey:UserID" json:"-"`
}

func (ReviewLike) TableName() string {
	return "review_likes"
}

type CommentLike struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CommentID uint64    `gorm:"not null" json:"comment_id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Comment ReviewComment `gorm:"foreignKey:CommentID" json:"-"`
	User    User          `gorm:"foreignKey:UserID" json:"-"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}