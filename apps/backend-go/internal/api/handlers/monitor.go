package handlers

import (
	"context"
	"net/http"
	"postpanda/backend-go/internal/database"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct{}

func NewMonitorHandler() *MonitorHandler { return &MonitorHandler{} }

func (h *MonitorHandler) Health(c *gin.Context) {
	ctx := context.Background()
	status := gin.H{"status": "ok", "services": gin.H{}}

	// Check DB
	if err := database.DB.Ping(ctx); err != nil {
		status["status"] = "degraded"
		status["services"].(gin.H)["database"] = "error"
	} else {
		status["services"].(gin.H)["database"] = "ok"
	}

	// Check Redis
	if err := database.Redis.Ping(ctx).Err(); err != nil {
		status["status"] = "degraded"
		status["services"].(gin.H)["redis"] = "error"
	} else {
		status["services"].(gin.H)["redis"] = "ok"
	}

	c.JSON(http.StatusOK, status)
}

func (h *MonitorHandler) GetQueueStats(c *gin.Context) {
	ctx := context.Background()
	queueName := c.Param("name")

	// Get queue length from Redis
	length, err := database.Redis.LLen(ctx, "queue:"+queueName).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"queue": queueName, "length": length})
}
