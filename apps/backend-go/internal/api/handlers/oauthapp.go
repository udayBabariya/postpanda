package handlers

import (
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type OAuthAppHandler struct {
	oauthAppSvc *services.OAuthAppService
}

func NewOAuthAppHandler(svc *services.OAuthAppService) *OAuthAppHandler {
	return &OAuthAppHandler{oauthAppSvc: svc}
}

func (h *OAuthAppHandler) Get(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	app, err := h.oauthAppSvc.GetByOrg(orgID)
	if err != nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, app)
}

func (h *OAuthAppHandler) Create(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		RedirectURL string `json:"redirectUrl" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	app, err := h.oauthAppSvc.Create(orgID, req.Name, req.Description, req.RedirectURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

func (h *OAuthAppHandler) Update(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		RedirectURL string `json:"redirectUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	app, err := h.oauthAppSvc.Update(orgID, id, req.Name, req.Description, req.RedirectURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

func (h *OAuthAppHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.oauthAppSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OAuth app deleted"})
}

// Public OAuth endpoints
func (h *OAuthAppHandler) AuthorizeOAuth(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")
	responseType := c.Query("response_type")

	if responseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid response_type"})
		return
	}

	code, err := h.oauthAppSvc.GenerateAuthCode(clientID, redirectURI, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, redirectURI+"?code="+code+"&state="+state)
}

func (h *OAuthAppHandler) TokenExchange(c *gin.Context) {
	var req struct {
		GrantType    string `json:"grant_type"`
		Code         string `json:"code"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	token, err := h.oauthAppSvc.ExchangeCode(req.Code, req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	})
}
