package models

import (
	"time"
)

type Cast struct {
	ID             uint64     `gorm:"primarykey" json:"id"`
	Name           string     `gorm:"size:255;not null" json:"name"`
	ProfileURL     *string    `gorm:"type:text" json:"profile_url,omitempty"`
	TmdbPersonID   *int       `gorm:"uniqueIndex" json:"tmdb_person_id,omitempty"`
	Bio            *string    `gorm:"type:text" json:"bio,omitempty"`
	BirthDate      *time.Time `gorm:"type:date" json:"birth_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	MovieCasts []MovieCast `gorm:"foreignKey:CastID" json:"-"`
	Movies     []Movie     `gorm:"many2many:movie_casts" json:"movies,omitempty"`
}

func (Cast) TableName() string {
	return "casts"
}

type MovieCast struct {
	ID            uint64    `gorm:"primarykey" json:"id"`
	MovieID       uint64    `gorm:"not null" json:"movie_id"`
	CastID        uint64    `gorm:"not null" json:"cast_id"`
	Role          string    `gorm:"size:100;not null" json:"role"` // 'actor', 'director', 'producer'
	CharacterName *string   `gorm:"size:255" json:"character_name,omitempty"`
	SortOrder     int       `gorm:"default:0" json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`

	// Relationships
	Movie Movie `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Cast  Cast  `gorm:"foreignKey:CastID" json:"cast,omitempty"`
}

func (MovieCast) TableName() string {
	return "movie_casts"
}
