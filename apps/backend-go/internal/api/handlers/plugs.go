package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type PlugsHandler struct {
	plugsSvc    *services.PlugsService
	customerSvc *services.CustomerService
}

func NewPlugsHandler(plugsSvc *services.PlugsService, customerSvc *services.CustomerService) *PlugsHandler {
	return &PlugsHandler{plugsSvc: plugsSvc, customerSvc: customerSvc}
}

func (h *PlugsHandler) ListAll(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	plugs, err := h.plugsSvc.ListByOrg(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plugs)
}

func (h *PlugsHandler) ListByIntegration(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")
	plugs, err := h.plugsSvc.ListByIntegration(orgID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plugs)
}

func (h *PlugsHandler) Upsert(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	integrationID := c.Param("id")

	var req struct {
		PlugFunction string `json:"plugFunction" binding:"required"`
		Data         string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	plug, err := h.plugsSvc.Upsert(orgID, integrationID, req.PlugFunction, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plug)
}

func (h *PlugsHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Activated bool `json:"activated"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.plugsSvc.Activate(id, req.Activated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plug updated"})
}

func (h *PlugsHandler) ListCustomers(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	customers, err := h.customerSvc.List(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *PlugsHandler) CreateCustomer(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	customer, err := h.customerSvc.Create(orgID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func (h *PlugsHandler) DeleteCustomer(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")
	if err := h.customerSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
