package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type SetsHandler struct {
	setSvc *services.SetService
}

func NewSetsHandler(setSvc *services.SetService) *SetsHandler {
	return &SetsHandler{setSvc: setSvc}
}

func (h *SetsHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	sets, err := h.setSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sets)
}

func (h *SetsHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Name    string `json:"name" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	set, err := h.setSvc.Create(orgID, req.Name, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, set)
}

func (h *SetsHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	set, err := h.setSvc.Update(orgID, id, req.Name, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, set)
}

func (h *SetsHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.setSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Set deleted"})
}
