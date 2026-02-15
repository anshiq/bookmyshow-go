package handler

import (
	"net/http"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) BookTicket(c *gin.Context) {
	var ticketDto models.TicketDto
	if err := c.ShouldBindJSON(&ticketDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ticket, err := h.ticketService.BookTicket(&ticketDto, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	tickets, err := h.ticketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) CancelTicket(c *gin.Context) {
	ticketID := c.Query("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	success, err := h.ticketService.CancelTicket(ticketID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}