package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminSvc *services.AdminService
}

func NewAdminHandler(svc *services.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: svc}
}

func (h *AdminHandler) GetErrors(c *gin.Context) {
	platform := c.Query("platform")
	page := c.DefaultQuery("page", "1")
	errors, total, err := h.adminSvc.GetErrors(platform, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": errors, "total": total})
}

func (h *AdminHandler) GetErrorPlatforms(c *gin.Context) {
	platforms, err := h.adminSvc.GetErrorPlatforms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, platforms)
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminSvc.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	search := c.Query("search")
	users, total, err := h.adminSvc.GetUsers(page, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "total": total})
}

func (h *AdminHandler) AddSubscription(c *gin.Context) {
	var req struct {
		OrgID    string `json:"orgId" binding:"required"`
		Tier     string `json:"tier" binding:"required"`
		Period   string `json:"period" binding:"required"`
		Channels int    `json:"channels"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.adminSvc.AddSubscription(req.OrgID, req.Tier, req.Period, req.Channels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription added"})
}

func (h *AdminHandler) CancelSubscription(c *gin.Context) {
	var req struct {
		OrgID string `json:"orgId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.adminSvc.CancelSubscription(req.OrgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled"})
}

func (h *AdminHandler) ImpersonateUser(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	token, err := h.adminSvc.ImpersonateUser(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
