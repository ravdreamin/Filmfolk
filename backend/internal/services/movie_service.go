package services

import (
	"errors"
	"fmt"

	"filmfolk/internal/db"
	"filmfolk/internal/models"

	"gorm.io/gorm"
)

// MovieService handles movie-related business logic
type MovieService struct{}

// NewMovieService creates a new movie service
func NewMovieService() *MovieService {
	return &MovieService{}
}



// GetMovie retrieves a movie by ID
func (s *MovieService) GetMovie(movieID uint64) (*models.Movie, error) {
	var movie models.Movie
	err := db.DB.First(&movie, movieID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &movie, nil
}

// ListMoviesFilter represents filter options for listing movies
type ListMoviesFilter struct {
	Genre    *string `form:"genre"`
	Year     *int    `form:"year"`
	Search   *string `form:"search"`
	SortBy   string  `form:"sort_by"` // rating, year, title, reviews
	Page     int     `form:"page"`
	PageSize int     `form:"page_size"`
}

// ListMovies retrieves movies with filters and pagination
func (s *MovieService) ListMovies(filter ListMoviesFilter) ([]models.Movie, int64, error) {
	// Default pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	query := db.DB.Model(&models.Movie{})

	// Apply filters
	if filter.Genre != nil && *filter.Genre != "" {
		query = query.Where("? = ANY(genres)", *filter.Genre)
	}

	if filter.Year != nil {
		query = query.Where("release_year = ?", *filter.Year)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("title ILIKE ?", searchTerm)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movies: %w", err)
	}

	// Apply sorting
	switch filter.SortBy {
	case "rating":
		query = query.Order("average_rating DESC NULLS LAST")
	case "year":
		query = query.Order("release_year DESC")
	case "reviews":
		query = query.Order("total_reviews DESC")
	default:
		query = query.Order("title ASC")
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// Execute query
	var movies []models.Movie
	if err := query.Find(&movies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch movies: %w", err)
	}

	return movies, total, nil
}

// UpdateMovieInput represents input for updating a movie
type UpdateMovieInput struct {
	Title          *string  `json:"title"`
	ReleaseYear    *int     `json:"release_year"`
	Genres         []string `json:"genres"`
	Summary        *string  `json:"summary"`
	PosterURL      *string  `json:"poster_url"`
	BackdropURL    *string  `json:"backdrop_url"`
	RuntimeMinutes *int     `json:"runtime_minutes"`
	Language       *string  `json:"language"`
}

// UpdateMovie updates movie information
func (s *MovieService) UpdateMovie(movieID uint64, input UpdateMovieInput) (*models.Movie, error) {
	var movie models.Movie
	if err := db.DB.First(&movie, movieID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.ReleaseYear != nil {
		updates["release_year"] = *input.ReleaseYear
	}
	if input.Genres != nil {
		updates["genres"] = input.Genres
	}
	if input.Summary != nil {
		updates["summary"] = *input.Summary
	}
	if input.PosterURL != nil {
		updates["poster_url"] = *input.PosterURL
	}
	if input.BackdropURL != nil {
		updates["backdrop_url"] = *input.BackdropURL
	}
	if input.RuntimeMinutes != nil {
		updates["runtime_minutes"] = *input.RuntimeMinutes
	}
	if input.Language != nil {
		updates["language"] = *input.Language
	}

	if err := db.DB.Model(&movie).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	// Reload movie
	if err := db.DB.First(&movie, movieID).Error; err != nil {
		return nil, err
	}

	return &movie, nil
}



// RecalculateMovieStats recalculates average rating and review count
func (s *MovieService) RecalculateMovieStats(movieID uint64) error {
	var stats struct {
		AvgRating    *float64
		ReviewCount  int64
	}

	err := db.DB.Model(&models.Review{}).
		Select("AVG(rating) as avg_rating, COUNT(*) as review_count").
		Where("movie_id = ?", movieID).
		Scan(&stats).Error

	if err != nil {
		return fmt.Errorf("failed to calculate stats: %w", err)
	}

	updates := map[string]interface{}{
		"total_reviews": stats.ReviewCount,
	}
	if stats.AvgRating != nil {
		updates["average_rating"] = *stats.AvgRating
	}

	return db.DB.Model(&models.Movie{}).Where("id = ?", movieID).Updates(updates).Error
}
