package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type IntegrationsHandler struct {
	integrationSvc *services.IntegrationService
}

func NewIntegrationsHandler(integrationSvc *services.IntegrationService) *IntegrationsHandler {
	return &IntegrationsHandler{integrationSvc: integrationSvc}
}

func (h *IntegrationsHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	integrations, err := h.integrationSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, integrations)
}

func (h *IntegrationsHandler) Get(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	integration, err := h.integrationSvc.GetByID(orgID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Integration not found"})
		return
	}

	c.JSON(http.StatusOK, integration)
}

func (h *IntegrationsHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.integrationSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Integration deleted"})
}

func (h *IntegrationsHandler) UpdatePostingTimes(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		PostingTimes string `json:"postingTimes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.integrationSvc.UpdatePostingTimes(orgID, id, req.PostingTimes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Posting times updated"})
}

func (h *IntegrationsHandler) Disable(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.integrationSvc.Disable(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Integration disabled"})
}

func (h *IntegrationsHandler) Enable(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.integrationSvc.Enable(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Integration enabled"})
}

func (h *IntegrationsHandler) GetProviders(c *gin.Context) {
	providers := h.integrationSvc.GetAvailableProviders()
	c.JSON(http.StatusOK, providers)
}

func (h *IntegrationsHandler) GetOAuthURL(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	provider := c.Param("provider")

	url, err := h.integrationSvc.GetOAuthURL(orgID, provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *IntegrationsHandler) OAuthCallback(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	integration, err := h.integrationSvc.HandleCallback(orgID, provider, code, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

func (h *IntegrationsHandler) RefreshToken(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.integrationSvc.RefreshToken(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

func (h *IntegrationsHandler) UpdateSettings(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.integrationSvc.UpdateSettings(orgID, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}
