package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsSvc *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsSvc *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsSvc: analyticsSvc}
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	data, err := h.analyticsSvc.GetDashboard(orgID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AnalyticsHandler) GetIntegrationStats(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	integrationID := c.Param("integrationId")

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	stats, err := h.analyticsSvc.GetIntegrationStats(orgID, integrationID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetPostStats(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("postId")

	stats, err := h.analyticsSvc.GetPostStats(orgID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetChannelOverview(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	overview, err := h.analyticsSvc.GetChannelOverview(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (h *AnalyticsHandler) GetScheduleAnalytics(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	data, err := h.analyticsSvc.GetScheduleAnalytics(orgID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
