package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type WebhooksHandler struct {
	webhookSvc *services.WebhookService
}

func NewWebhooksHandler(webhookSvc *services.WebhookService) *WebhooksHandler {
	return &WebhooksHandler{webhookSvc: webhookSvc}
}

func (h *WebhooksHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	webhooks, err := h.webhookSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhooks)
}

func (h *WebhooksHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Name string   `json:"name" binding:"required"`
		URL  string   `json:"url" binding:"required,url"`
		Integrations []string `json:"integrations"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	webhook, err := h.webhookSvc.Create(orgID, req.Name, req.URL, req.Integrations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, webhook)
}

func (h *WebhooksHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		Name string   `json:"name"`
		URL  string   `json:"url"`
		Integrations []string `json:"integrations"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	webhook, err := h.webhookSvc.Update(orgID, id, req.Name, req.URL, req.Integrations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

func (h *WebhooksHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.webhookSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook deleted"})
}
