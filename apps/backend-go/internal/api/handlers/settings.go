package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsSvc *services.SettingsService
}

func NewSettingsHandler(settingsSvc *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsSvc: settingsSvc}
}

func (h *SettingsHandler) GetOrgSettings(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	settings, err := h.settingsSvc.GetOrgSettings(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) UpdateOrgSettings(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.settingsSvc.UpdateOrgSettings(orgID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}

func (h *SettingsHandler) RegenerateAPIKey(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	apiKey, err := h.settingsSvc.RegenerateAPIKey(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"apiKey": apiKey})
}

func (h *SettingsHandler) UpdateShortLink(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		ShortLink string `json:"shortlink" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.settingsSvc.UpdateShortLink(orgID, req.ShortLink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Short link preference updated"})
}
