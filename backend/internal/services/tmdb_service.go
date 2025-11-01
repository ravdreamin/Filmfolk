package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// TMDBService handles communication with The Movie Database API
type TMDBService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewTMDBService creates a new TMDB service
func NewTMDBService(apiKey string) *TMDBService {
	return &TMDBService{
		apiKey:  apiKey,
		baseURL: "https://api.themoviedb.org/3",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// TMDBMovie represents a movie from TMDB API
type TMDBMovie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Overview    string   `json:"overview"`
	ReleaseDate string   `json:"release_date"`
	PosterPath  string   `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	VoteAverage float64  `json:"vote_average"`
	Genres      []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Runtime         int    `json:"runtime"`
	OriginalLanguage string `json:"original_language"`
	IMDbID          string `json:"imdb_id"`
}

// TMDBSearchResult represents search results from TMDB
type TMDBSearchResult struct {
	Page    int `json:"page"`
	Results []struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Overview    string  `json:"overview"`
		ReleaseDate string  `json:"release_date"`
		PosterPath  string  `json:"poster_path"`
		VoteAverage float64 `json:"vote_average"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// TMDBCast represents cast member from TMDB
type TMDBCast struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
	Order       int    `json:"order"`
}

// TMDBCrew represents crew member from TMDB
type TMDBCrew struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
}

// TMDBCredits represents credits (cast and crew) from TMDB
type TMDBCredits struct {
	ID   int        `json:"id"`
	Cast []TMDBCast `json:"cast"`
	Crew []TMDBCrew `json:"crew"`
}

// SearchMovies searches for movies on TMDB
func (s *TMDBService) SearchMovies(query string, page int) (*TMDBSearchResult, error) {
	if s.apiKey == "" {
		return nil, errors.New("TMDB API key not configured")
	}

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("query", query)
	params.Add("page", fmt.Sprintf("%d", page))

	url := fmt.Sprintf("%s/search/movie?%s", s.baseURL, params.Encode())

	var result TMDBSearchResult
	if err := s.makeRequest(url, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMovieDetails fetches detailed information about a movie
func (s *TMDBService) GetMovieDetails(tmdbID int) (*TMDBMovie, error) {
	if s.apiKey == "" {
		return nil, errors.New("TMDB API key not configured")
	}

	url := fmt.Sprintf("%s/movie/%d?api_key=%s", s.baseURL, tmdbID, s.apiKey)

	var movie TMDBMovie
	if err := s.makeRequest(url, &movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

// GetMovieCredits fetches cast and crew for a movie
func (s *TMDBService) GetMovieCredits(tmdbID int) (*TMDBCredits, error) {
	if s.apiKey == "" {
		return nil, errors.New("TMDB API key not configured")
	}

	url := fmt.Sprintf("%s/movie/%d/credits?api_key=%s", s.baseURL, tmdbID, s.apiKey)

	var credits TMDBCredits
	if err := s.makeRequest(url, &credits); err != nil {
		return nil, err
	}

	return &credits, nil
}

// GetImageURL constructs full URL for TMDB images
func (s *TMDBService) GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	// size can be: w92, w154, w185, w342, w500, w780, original
	return fmt.Sprintf("https://image.tmdb.org/t/p/%s%s", size, path)
}

// makeRequest is a helper to make HTTP requests to TMDB
func (s *TMDBService) makeRequest(url string, result interface{}) error {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("TMDB API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TMDB API error (status %d): %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode TMDB response: %w", err)
	}

	return nil
}
