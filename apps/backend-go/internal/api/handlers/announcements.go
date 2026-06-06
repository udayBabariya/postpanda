package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type AnnouncementsHandler struct {
	announcementSvc *services.AnnouncementService
}

func NewAnnouncementsHandler(svc *services.AnnouncementService) *AnnouncementsHandler {
	return &AnnouncementsHandler{announcementSvc: svc}
}

func (h *AnnouncementsHandler) List(c *gin.Context) {
	announcements, err := h.announcementSvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, announcements)
}

func (h *AnnouncementsHandler) Create(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "INFO"
	}

	announcement, err := h.announcementSvc.Create(req.Title, req.Description, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, announcement)
}

func (h *AnnouncementsHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.announcementSvc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement deleted"})
}
