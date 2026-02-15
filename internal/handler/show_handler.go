package handler

import (
	"net/http"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

type ShowHandler struct {
	showService service.ShowService
}

func NewShowHandler(showService service.ShowService) *ShowHandler {
	return &ShowHandler{showService: showService}
}

func (h *ShowHandler) AddShow(c *gin.Context) {
	var showDto models.ShowDto
	if err := c.ShouldBindJSON(&showDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	show, err := h.showService.AddShow(&showDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, show)
}

func (h *ShowHandler) GetAllShows(c *gin.Context) {
	shows, err := h.showService.GetAllShows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

func (h *ShowHandler) GetMovieShows(c *gin.Context) {
	movieID := c.Query("id")
	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter required"})
		return
	}

	shows, err := h.showService.GetMovieShows(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

func (h *ShowHandler) GetShowByID(c *gin.Context) {
	id := c.Param("id")
	show, err := h.showService.GetShowByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	c.JSON(http.StatusOK, show)
}