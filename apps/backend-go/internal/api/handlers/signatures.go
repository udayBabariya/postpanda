package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type SignaturesHandler struct {
	sigSvc *services.SignatureService
}

func NewSignaturesHandler(sigSvc *services.SignatureService) *SignaturesHandler {
	return &SignaturesHandler{sigSvc: sigSvc}
}

func (h *SignaturesHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	sigs, err := h.sigSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sigs)
}

func (h *SignaturesHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Content string `json:"content" binding:"required"`
		AutoAdd bool   `json:"autoAdd"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sig, err := h.sigSvc.Create(orgID, req.Content, req.AutoAdd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sig)
}

func (h *SignaturesHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		Content string `json:"content"`
		AutoAdd bool   `json:"autoAdd"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sig, err := h.sigSvc.Update(orgID, id, req.Content, req.AutoAdd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sig)
}

func (h *SignaturesHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.sigSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signature deleted"})
}
