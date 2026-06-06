package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type PublicAPIHandler struct {
	postSvc    *services.PostService
	commentSvc *services.CommentService
}

func NewPublicAPIHandler(postSvc *services.PostService, commentSvc *services.CommentService) *PublicAPIHandler {
	return &PublicAPIHandler{postSvc: postSvc, commentSvc: commentSvc}
}

func (h *PublicAPIHandler) GetPostPreview(c *gin.Context) {
	postID := c.Param("id")
	// Public preview - no auth needed, return limited fields
	post, err := h.postSvc.GetPublicPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PublicAPIHandler) GetPostComments(c *gin.Context) {
	postID := c.Param("id")
	comments, err := h.commentSvc.ListByPost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *PublicAPIHandler) CreateComment(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)
	postID := c.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	comment, err := h.commentSvc.Create(orgID, userID, postID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *PublicAPIHandler) CanRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"canRegister": !config.App.DisableRegistration})
}
