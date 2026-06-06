package handlers

import (
	"io"
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81/webhook"
)

type BillingHandler struct {
	billingSvc *services.BillingService
}

func NewBillingHandler(billingSvc *services.BillingService) *BillingHandler {
	return &BillingHandler{billingSvc: billingSvc}
}

func (h *BillingHandler) GetSubscription(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	sub, err := h.billingSvc.GetSubscription(orgID)
	if err != nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *BillingHandler) CreateCheckout(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	var req struct {
		Period string `json:"period" binding:"required"`
		Plan   string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	url, err := h.billingSvc.CreateCheckoutSession(orgID, userID, req.Plan, req.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *BillingHandler) CreatePortal(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	url, err := h.billingSvc.CreatePortalSession(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *BillingHandler) StripeWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read body"})
		return
	}

	sig := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sig, config.App.StripeWebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signature"})
		return
	}

	if err := h.billingSvc.HandleWebhookEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	if err := h.billingSvc.CancelSubscription(orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled"})
}

func (h *BillingHandler) GetPlans(c *gin.Context) {
	plans := h.billingSvc.GetPlans()
	c.JSON(http.StatusOK, plans)
}

func (h *BillingHandler) CheckDiscount(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	result, err := h.billingSvc.GetDiscountEligibility(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BillingHandler) ApplyDiscount(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.billingSvc.ApplyDiscount(orgID, req.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discount applied"})
}

func (h *BillingHandler) FinishTrial(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	if err := h.billingSvc.FinishTrial(orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trial finished"})
}

func (h *BillingHandler) IsTrialFinished(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	finished, err := h.billingSvc.IsTrialFinished(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"finished": finished})
}

func (h *BillingHandler) ApplyLifetimeDeal(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.billingSvc.ApplyLifetimeDeal(orgID, req.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lifetime deal applied"})
}
