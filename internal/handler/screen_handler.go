package handler

import (
	"net/http"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

type ScreenHandler struct {
	screenService service.ScreenService
}

func NewScreenHandler(screenService service.ScreenService) *ScreenHandler {
	return &ScreenHandler{screenService: screenService}
}

func (h *ScreenHandler) AddScreen(c *gin.Context) {
	var screenDto models.ScreenDto
	if err := c.ShouldBindJSON(&screenDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	screen, err := h.screenService.AddScreen(&screenDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, screen)
}

func (h *ScreenHandler) GetAllScreens(c *gin.Context) {
	screens, err := h.screenService.GetAllScreens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, screens)
}

func (h *ScreenHandler) GetScreensByTheatreID(c *gin.Context) {
	theatreID := c.Param("id")
	screens, err := h.screenService.GetScreensByTheatreID(theatreID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, screens)
}