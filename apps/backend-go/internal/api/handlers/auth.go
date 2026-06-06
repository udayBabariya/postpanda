package handlers

import (
	"context"
	"net/http"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/auth"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/database"
	"postpanda/backend-go/internal/models"
	"postpanda/backend-go/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userSvc *services.UserService
}

func NewAuthHandler(userSvc *services.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	if config.App.DisableRegistration {
		c.JSON(http.StatusForbidden, gin.H{"message": "Registration is disabled"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	existing, _ := h.userSvc.FindByEmail(req.Email, models.ProviderLocal)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already in use"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	user, org, err := h.userSvc.CreateWithOrganization(req.Email, string(hashed), req.Name, models.ProviderLocal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	token, err := auth.GenerateToken(user.ID, org.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusCreated, gin.H{
		"token":        token,
		"userId":       user.ID,
		"orgId":        org.ID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.userSvc.FindByEmail(req.Email, models.ProviderLocal)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	if user.Password == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	if !user.Activated {
		c.JSON(http.StatusForbidden, gin.H{"message": "Account not activated"})
		return
	}

	org, err := h.userSvc.GetPrimaryOrg(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get organization"})
		return
	}

	token, err := auth.GenerateToken(user.ID, org.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{
		"token":  token,
		"userId": user.ID,
		"orgId":  org.ID,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("postiz-auth", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orgID := middleware.GetOrgID(c)

	user, err := h.userSvc.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	org, err := h.userSvc.FindOrgByID(orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"org":  org,
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.userSvc.SendPasswordReset(req.Email); err != nil {
		// Always return OK to avoid email enumeration
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.userSvc.ResetPassword(req.Token, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *AuthHandler) ActivateAccount(c *gin.Context) {
	code := c.Param("code")
	if err := h.userSvc.ActivateAccount(code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid activation code"})
		return
	}
	c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?activated=true")
}

func (h *AuthHandler) SwitchOrg(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orgID := c.Param("orgId")

	membership, err := h.userSvc.GetOrgMembership(userID, orgID)
	if err != nil || membership == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Access denied"})
		return
	}

	token, err := auth.GenerateToken(userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// OAuth providers

func (h *AuthHandler) GitHubOAuth(c *gin.Context) {
	provider := auth.NewGitHubProvider()
	url := provider.GetAuthURL(generateStateToken())
	c.Redirect(http.StatusFound, url)
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?error=github_denied")
		return
	}

	provider := auth.NewGitHubProvider()
	userInfo, err := provider.ExchangeCode(code)
	if err != nil {
		c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?error=github_failed")
		return
	}

	user, org, err := h.userSvc.FindOrCreateOAuthUser(userInfo.Email, userInfo.Name, userInfo.ID, models.ProviderGitHub)
	if err != nil {
		c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?error=user_failed")
		return
	}

	token, _ := auth.GenerateToken(user.ID, org.ID)
	setAuthCookie(c, token)
	c.Redirect(http.StatusFound, config.App.FrontendURL+"?token="+token)
}

func (h *AuthHandler) GoogleOAuth(c *gin.Context) {
	provider := auth.NewGoogleProvider()
	url := provider.GetAuthURL(generateStateToken())
	c.Redirect(http.StatusFound, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?error=google_denied")
		return
	}

	provider := auth.NewGoogleProvider()
	userInfo, err := provider.ExchangeCode(code)
	if err != nil {
		c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?error=google_failed")
		return
	}

	user, org, err := h.userSvc.FindOrCreateOAuthUser(userInfo.Email, userInfo.Name, userInfo.ID, models.ProviderGoogle)
	if err != nil {
		c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth?error=user_failed")
		return
	}

	token, _ := auth.GenerateToken(user.ID, org.ID)
	setAuthCookie(c, token)
	c.Redirect(http.StatusFound, config.App.FrontendURL+"?token="+token)
}

// Update last online timestamp
func (h *AuthHandler) UpdateLastOnline(userID string) {
	ctx := context.Background()
	database.DB.Exec(ctx, "UPDATE users SET last_online = $1 WHERE id = $2", time.Now(), userID)
}

func setAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("postiz-auth", token, 30*24*3600, "/", "", false, true)
}

func generateStateToken() string {
	return uuid.New().String()
}

func (h *AuthHandler) CanRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"canRegister": !config.App.DisableRegistration,
	})
}

func (h *AuthHandler) ResendActivation(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Resend activation - always return OK to prevent email enumeration
	h.userSvc.ResendActivation(req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "If the email exists and is unactivated, a new activation email has been sent"})
}

func (h *AuthHandler) CheckOAuthExists(c *gin.Context) {
	provider := c.Param("provider")
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	exists := h.userSvc.CheckOAuthExists(req.ID, models.Provider(provider))
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (h *AuthHandler) OAuthMobileCallback(c *gin.Context) {
	provider := c.Query("provider")
	code := c.Query("code")
	_ = provider
	_ = code
	// Handle mobile OAuth deep link callback
	c.Redirect(http.StatusFound, config.App.FrontendURL+"/auth/callback?code="+code)
}
