package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/auth"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct {
	userSvc    *services.UserService
	billingSvc *services.BillingService
}

func NewEnterpriseHandler(userSvc *services.UserService, billingSvc *services.BillingService) *EnterpriseHandler {
	return &EnterpriseHandler{userSvc: userSvc, billingSvc: billingSvc}
}

func (h *EnterpriseHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		OrgID    string `json:"orgId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, org, err := h.userSvc.CreateWithOrganization(req.Email, req.Password, req.Name, "LOCAL")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	token, _ := auth.GenerateToken(user.ID, org.ID)
	c.JSON(http.StatusCreated, gin.H{"token": token, "userId": user.ID, "orgId": org.ID})
}

func (h *EnterpriseHandler) GetRedirectParams(c *gin.Context) {
	var req struct {
		OrgID string `json:"orgId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orgId": req.OrgID})
}

func (h *EnterpriseHandler) DeleteChannel(c *gin.Context) {
	var req struct {
		IntegrationID string `json:"integrationId" binding:"required"`
		OrgID         string `json:"orgId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	integrationSvc := services.NewIntegrationService()
	if err := integrationSvc.Delete(req.OrgID, req.IntegrationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted"})
}
