package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type CopilotHandler struct {
	copilotSvc *services.CopilotService
}

func NewCopilotHandler(svc *services.CopilotService) *CopilotHandler {
	return &CopilotHandler{copilotSvc: svc}
}

func (h *CopilotHandler) Chat(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	var req struct {
		Message  string `json:"message" binding:"required"`
		ThreadID string `json:"threadId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resp, err := h.copilotSvc.Chat(orgID, userID, req.Message, req.ThreadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *CopilotHandler) GetCredits(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	credits, err := h.copilotSvc.GetCredits(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"credits": credits})
}

func (h *CopilotHandler) ListThreads(c *gin.Context) {
	userID := middleware.GetUserID(c)
	threads, err := h.copilotSvc.ListThreads(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, threads)
}

func (h *CopilotHandler) GetThreadMessages(c *gin.Context) {
	threadID := c.Param("threadId")
	messages, err := h.copilotSvc.GetThreadMessages(threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *CopilotHandler) GenerateDraft(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Topic       string   `json:"topic" binding:"required"`
		Platform    string   `json:"platform"`
		Tone        string   `json:"tone"`
		Integration string   `json:"integration"`
		Keywords    []string `json:"keywords"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	drafts, err := h.copilotSvc.GenerateDraft(orgID, req.Topic, req.Platform, req.Tone, req.Keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"drafts": drafts})
}

func (h *CopilotHandler) GenerateImage(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Prompt string `json:"prompt" binding:"required"`
		Style  string `json:"style"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.copilotSvc.GenerateImage(orgID, req.Prompt, req.Style)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *CopilotHandler) SeparatePosts(c *gin.Context) {
	var req struct {
		Content  string `json:"content" binding:"required"`
		MaxChars int    `json:"maxChars"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if req.MaxChars == 0 {
		req.MaxChars = 280
	}
	parts := h.copilotSvc.SeparatePosts(req.Content, req.MaxChars)
	c.JSON(http.StatusOK, gin.H{"parts": parts})
}
