package handler

import (
	"net/http"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configService service.ConfigService
}

func NewConfigHandler(configService service.ConfigService) *ConfigHandler {
	return &ConfigHandler{configService: configService}
}

func (h *ConfigHandler) GetByDBType(c *gin.Context) {
	dbType := c.Query("dbType")
	if dbType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dbType parameter required"})
		return
	}

	configs, err := h.configService.GetByDBType(dbType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ConfigHandler) GetByDBTypeAndHashID(c *gin.Context) {
	dbType := c.Query("dbType")
	hashID := c.Query("hashId")

	if dbType == "" || hashID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dbType and hashId parameters required"})
		return
	}

	config, err := h.configService.GetByDBTypeAndHashID(dbType, hashID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ConfigHandler) SetConfig(c *gin.Context) {
	var configDto models.ConfigDto
	if err := c.ShouldBindJSON(&configDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.configService.SetConfig(&configDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}
