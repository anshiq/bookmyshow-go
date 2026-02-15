package handler

import (
	"net/http"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieService service.MovieService
}

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

func (h *MovieHandler) AddMovie(c *gin.Context) {
	var movieDto models.MovieDto
	if err := c.ShouldBindJSON(&movieDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := h.movieService.AddMovie(&movieDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.movieService.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(movies) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No movies found"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) SearchMovies(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter required"})
		return
	}

	movies, err := h.movieService.SearchMovies(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id := c.Param("id")
	movie, err := h.movieService.GetMovieByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}