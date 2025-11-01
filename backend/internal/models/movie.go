package models

import (
	"time"

	"github.com/lib/pq"
)

type MovieStatus string

const (
	MovieStatusPending  MovieStatus = "pending_approval"
	MovieStatusApproved MovieStatus = "approved"
	MovieStatusRejected MovieStatus = "rejected"
)

type Movie struct {
	ID             uint64         `gorm:"primarykey" json:"id"`
	Title          string         `gorm:"size:500;not null" json:"title"`
	ReleaseYear    int            `gorm:"not null" json:"release_year"`
	Genres         pq.StringArray `gorm:"type:varchar(255)[]" json:"genres"`
	Summary        *string        `gorm:"type:text" json:"summary,omitempty"`
	PosterURL      *string        `gorm:"type:text" json:"poster_url,omitempty"`
	BackdropURL    *string        `gorm:"type:text" json:"backdrop_url,omitempty"`
	RuntimeMinutes *int           `json:"runtime_minutes,omitempty"`
	Language       *string        `gorm:"size:50" json:"language,omitempty"`

	// External API integration
	TmdbID *int    `gorm:"uniqueIndex" json:"tmdb_id,omitempty"`
	ImdbID *string `gorm:"size:20" json:"imdb_id,omitempty"`

	Status            MovieStatus `gorm:"type:movie_status;not null;default:pending_approval" json:"status"`
	SubmittedByUserID *uint64     `json:"submitted_by_user_id,omitempty"`
	ApprovedByUserID  *uint64     `json:"approved_by_user_id,omitempty"`

	// Aggregated stats
	AverageRating *float64 `gorm:"type:decimal(3,2)" json:"average_rating,omitempty"`
	TotalReviews  int      `gorm:"default:0" json:"total_reviews"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	SubmittedBy *User    `gorm:"foreignKey:SubmittedByUserID" json:"submitted_by,omitempty"`
	ApprovedBy  *User    `gorm:"foreignKey:ApprovedByUserID" json:"approved_by,omitempty"`
	Reviews     []Review `gorm:"foreignKey:MovieID" json:"-"`
}

func (Movie) TableName() string {
	return "movies"
}
