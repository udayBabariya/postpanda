package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type AutoPostHandler struct {
	autoPostSvc *services.AutoPostService
}

func NewAutoPostHandler(autoPostSvc *services.AutoPostService) *AutoPostHandler {
	return &AutoPostHandler{autoPostSvc: autoPostSvc}
}

func (h *AutoPostHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	autoPosts, err := h.autoPostSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, autoPosts)
}

func (h *AutoPostHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Title           string   `json:"title" binding:"required"`
		Content         string   `json:"content"`
		OnSlot          bool     `json:"onSlot"`
		SyncLast        bool     `json:"syncLast"`
		URL             string   `json:"url" binding:"required"`
		Active          bool     `json:"active"`
		AddPicture      bool     `json:"addPicture"`
		GenerateContent bool     `json:"generateContent"`
		Integrations    []string `json:"integrations"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	autoPost, err := h.autoPostSvc.Create(orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, autoPost)
}

func (h *AutoPostHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	autoPost, err := h.autoPostSvc.Update(orgID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, autoPost)
}

func (h *AutoPostHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.autoPostSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto-post deleted"})
}

func (h *AutoPostHandler) Toggle(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.autoPostSvc.Toggle(orgID, id, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto-post toggled"})
}
