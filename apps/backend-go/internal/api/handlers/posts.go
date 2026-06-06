package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type PostsHandler struct {
	postSvc *services.PostService
}

func NewPostsHandler(postSvc *services.PostService) *PostsHandler {
	return &PostsHandler{postSvc: postSvc}
}

type CreatePostRequest struct {
	Date         string        `json:"date" binding:"required"`
	Integrations []string      `json:"integrations" binding:"required"`
	Posts        []PostContent `json:"posts" binding:"required"`
	Tags         []string      `json:"tags"`
	Settings     interface{}   `json:"settings"`
}

type PostContent struct {
	Content  string        `json:"content"`
	Image    []interface{} `json:"image"`
	Settings interface{}   `json:"settings"`
}

func (h *PostsHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	posts, err := h.postSvc.Create(orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, posts)
}

func (h *PostsHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	integrationID := c.Query("integrationId")

	posts, err := h.postSvc.List(orgID, startDate, endDate, integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostsHandler) Get(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	post, err := h.postSvc.GetByID(orgID, postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostsHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	post, err := h.postSvc.Update(orgID, postID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostsHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	if err := h.postSvc.Delete(orgID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func (h *PostsHandler) GetGroup(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	groupID := c.Param("groupId")

	posts, err := h.postSvc.GetGroup(orgID, groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post group not found"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostsHandler) DeleteGroup(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	groupID := c.Param("groupId")

	if err := h.postSvc.DeleteGroup(orgID, groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

func (h *PostsHandler) Reschedule(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	var req struct {
		Date string `json:"date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	post, err := h.postSvc.Reschedule(orgID, postID, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostsHandler) SubmitForApproval(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	var req struct {
		OrderID string `json:"orderId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.postSvc.SubmitForApproval(orgID, postID, req.OrderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submitted for approval"})
}

func (h *PostsHandler) Approve(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	if err := h.postSvc.Approve(orgID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post approved"})
}

func (h *PostsHandler) GetAnalytics(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	postID := c.Param("id")

	analytics, err := h.postSvc.GetAnalytics(orgID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
