package models

import (
	"time"
)

type MovieStatus string

const (
	MovieStatusPendingApproval MovieStatus = "pending_approval"
	MovieStatusApproved        MovieStatus = "approved"
	MovieStatusRejected        MovieStatus = "rejected"
)

type Movie struct {
	ID            uint64      `gorm:"primaryKey"`
	Title         string      `gorm:"uniqueIndex:idx_title_year;not null"`
	ReleaseYear   int         `gorm:"uniqueIndex:idx_title_year;not null"`
	Genre         string      `gorm:"column:genre"`
	Summary       string      `gorm:"column:summary"`
	ExternalApiID string      `gorm:"column:external_api_id;unique"`
	Status        MovieStatus `gorm:"type:movie_status;not null;default:'pending_approval'"`
	CreatedAt     time.Time   `gorm:"column:created_at"`
	UpdatedAt     time.Time   `gorm:"column:updated_at"`
}