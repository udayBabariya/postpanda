package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	userSvc *services.UserService
}

func NewUsersHandler(userSvc *services.UserService) *UsersHandler {
	return &UsersHandler{userSvc: userSvc}
}

func (h *UsersHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userSvc.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UsersHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Name      string `json:"name"`
		LastName  string `json:"lastName"`
		Bio       string `json:"bio"`
		Timezone  int    `json:"timezone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"name":     req.Name,
		"lastName": req.LastName,
		"bio":      req.Bio,
		"timezone": req.Timezone,
	}

	if err := h.userSvc.UpdateProfile(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, _ := h.userSvc.FindByID(userID)
	c.JSON(http.StatusOK, user)
}

func (h *UsersHandler) UpdateEmailPreferences(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		SendSuccessEmails bool `json:"sendSuccessEmails"`
		SendFailureEmails bool `json:"sendFailureEmails"`
		SendStreakEmails  bool `json:"sendStreakEmails"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.userSvc.UpdateEmailPreferences(userID, req.SendSuccessEmails, req.SendFailureEmails, req.SendStreakEmails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email preferences updated"})
}

func (h *UsersHandler) GetOrganizations(c *gin.Context) {
	userID := middleware.GetUserID(c)

	orgs, err := h.userSvc.ListOrgs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

func (h *UsersHandler) GetOrgMembers(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	members, err := h.userSvc.GetOrgMembers(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *UsersHandler) InviteMember(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.userSvc.InviteMember(orgID, req.Email, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent"})
}

func (h *UsersHandler) RemoveMember(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	memberID := c.Param("memberId")

	if err := h.userSvc.RemoveMember(orgID, memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

func (h *UsersHandler) UpdateMemberRole(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	memberID := c.Param("memberId")

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.userSvc.UpdateMemberRole(orgID, memberID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h *UsersHandler) CreateOrg(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	org, err := h.userSvc.CreateOrg(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}
