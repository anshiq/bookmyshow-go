package handler

import (
	"net/http"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

type TheatreHandler struct {
	theatreService service.TheatreService
}

func NewTheatreHandler(theatreService service.TheatreService) *TheatreHandler {
	return &TheatreHandler{theatreService: theatreService}
}

func (h *TheatreHandler) AddTheatre(c *gin.Context) {
	var theatreDto models.TheatreDto
	if err := c.ShouldBindJSON(&theatreDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	theatre, err := h.theatreService.AddTheatre(&theatreDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, theatre)
}

func (h *TheatreHandler) GetAllTheatres(c *gin.Context) {
	theatres, err := h.theatreService.GetAllTheatres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, theatres)
}

func (h *TheatreHandler) GetTheatreByID(c *gin.Context) {
	id := c.Param("id")
	theatre, err := h.theatreService.GetTheatreByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theatre not found"})
		return
	}

	c.JSON(http.StatusOK, theatre)
}