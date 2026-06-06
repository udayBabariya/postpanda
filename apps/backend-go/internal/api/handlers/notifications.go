package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationsHandler struct {
	notifSvc *services.NotificationService
}

func NewNotificationsHandler(notifSvc *services.NotificationService) *NotificationsHandler {
	return &NotificationsHandler{notifSvc: notifSvc}
}

func (h *NotificationsHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	notifications, err := h.notifSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationsHandler) MarkRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.notifSvc.MarkRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications marked as read"})
}

func (h *NotificationsHandler) GetUnreadCount(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	count, err := h.notifSvc.GetUnreadCount(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
