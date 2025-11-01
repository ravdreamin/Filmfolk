package handlers

import (
	"net/http"
	"strconv"

	"filmfolk/internal/services"

	"github.com/gin-gonic/gin"
)

// MovieHandler handles movie-related HTTP requests
type MovieHandler struct {
	movieService *services.MovieService
}

// NewMovieHandler creates a new movie handler
func NewMovieHandler() *MovieHandler {
	return &MovieHandler{
		movieService: services.NewMovieService(),
	}
}

// GetMovie handles GET /api/v1/movies/:id
func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := h.movieService.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// ListMovies handles GET /api/v1/movies
func (h *MovieHandler) ListMovies(c *gin.Context) {
	var filter services.ListMoviesFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movies, total, err := h.movieService.ListMovies(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"total":  total,
		"page":   filter.Page,
		"page_size": filter.PageSize,
	})
}

// UpdateMovie handles PUT /api/v1/movies/:id
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var input services.UpdateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := h.movieService.UpdateMovie(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}