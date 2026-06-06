package api

import (
	"postpanda/backend-go/internal/api/handlers"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	return &Server{router: r}
}

func (s *Server) SetupRoutes() {
	// Services
	userSvc := services.NewUserService()
	postSvc := services.NewPostService()
	integrationSvc := services.NewIntegrationService()
	mediaSvc := services.NewMediaService()
	analyticsSvc := services.NewAnalyticsService()
	tagSvc := services.NewTagService()
	notifSvc := services.NewNotificationService()
	webhookSvc := services.NewWebhookService()
	billingSvc := services.NewBillingService()
	settingsSvc := services.NewSettingsService()
	setSvc := services.NewSetService()
	sigSvc := services.NewSignatureService()
	autoPostSvc := services.NewAutoPostService()
	announcementSvc := services.NewAnnouncementService()
	oauthAppSvc := services.NewOAuthAppService()

	// Handlers
	authH := handlers.NewAuthHandler(userSvc)
	postsH := handlers.NewPostsHandler(postSvc)
	integrationsH := handlers.NewIntegrationsHandler(integrationSvc)
	mediaH := handlers.NewMediaHandler(mediaSvc)
	analyticsH := handlers.NewAnalyticsHandler(analyticsSvc)
	usersH := handlers.NewUsersHandler(userSvc)
	tagsH := handlers.NewTagsHandler(tagSvc)
	notifH := handlers.NewNotificationsHandler(notifSvc)
	webhooksH := handlers.NewWebhooksHandler(webhookSvc)
	billingH := handlers.NewBillingHandler(billingSvc)
	settingsH := handlers.NewSettingsHandler(settingsSvc)
	setsH := handlers.NewSetsHandler(setSvc)
	sigsH := handlers.NewSignaturesHandler(sigSvc)
	autoPostH := handlers.NewAutoPostHandler(autoPostSvc)
	announcementsH := handlers.NewAnnouncementsHandler(announcementSvc)
	oauthAppH := handlers.NewOAuthAppHandler(oauthAppSvc)

	// Middleware
	s.router.Use(middleware.CORS())
	s.router.Use(middleware.RateLimit())

	// Static files
	s.router.Static("/uploads", "./uploads")

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (no auth required)
	authGroup := s.router.Group("/auth")
	{
		authGroup.POST("/register", authH.Register)
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/logout", authH.Logout)
		authGroup.POST("/forgot-password", authH.ForgotPassword)
		authGroup.POST("/reset-password", authH.ResetPassword)
		authGroup.GET("/activate/:code", authH.ActivateAccount)

		// OAuth providers
		authGroup.GET("/github", authH.GitHubOAuth)
		authGroup.GET("/github/callback", authH.GitHubCallback)
		authGroup.GET("/google", authH.GoogleOAuth)
		authGroup.GET("/google/callback", authH.GoogleCallback)
	}

	// Stripe webhook (raw body needed)
	s.router.POST("/stripe/webhook", billingH.StripeWebhook)

	// Public OAuth app endpoints
	s.router.GET("/oauth/authorize", oauthAppH.AuthorizeOAuth)
	s.router.POST("/oauth/token", oauthAppH.TokenExchange)

	// Authenticated API routes
	api := s.router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		// Auth
		api.GET("/auth/me", authH.Me)
		api.POST("/auth/switch/:orgId", authH.SwitchOrg)

		// Users
		api.GET("/users/profile", usersH.GetProfile)
		api.PUT("/users/profile", usersH.UpdateProfile)
		api.PUT("/users/email-preferences", usersH.UpdateEmailPreferences)
		api.GET("/users/organizations", usersH.GetOrganizations)
		api.POST("/users/organizations", usersH.CreateOrg)
		api.GET("/users/organizations/members", usersH.GetOrgMembers)
		api.POST("/users/organizations/members/invite", usersH.InviteMember)
		api.DELETE("/users/organizations/members/:memberId", usersH.RemoveMember)
		api.PATCH("/users/organizations/members/:memberId/role", usersH.UpdateMemberRole)

		// Posts
		api.GET("/posts", postsH.List)
		api.POST("/posts", postsH.Create)
		api.GET("/posts/:id", postsH.Get)
		api.PUT("/posts/:id", postsH.Update)
		api.DELETE("/posts/:id", postsH.Delete)
		api.POST("/posts/:id/reschedule", postsH.Reschedule)
		api.POST("/posts/:id/submit", postsH.SubmitForApproval)
		api.POST("/posts/:id/approve", postsH.Approve)
		api.GET("/posts/:id/analytics", postsH.GetAnalytics)
		api.GET("/posts/group/:groupId", postsH.GetGroup)
		api.DELETE("/posts/group/:groupId", postsH.DeleteGroup)

		// Integrations
		api.GET("/integrations", integrationsH.List)
		api.GET("/integrations/providers", integrationsH.GetProviders)
		api.GET("/integrations/:id", integrationsH.Get)
		api.DELETE("/integrations/:id", integrationsH.Delete)
		api.PATCH("/integrations/:id/posting-times", integrationsH.UpdatePostingTimes)
		api.PATCH("/integrations/:id/disable", integrationsH.Disable)
		api.PATCH("/integrations/:id/enable", integrationsH.Enable)
		api.POST("/integrations/:id/refresh", integrationsH.RefreshToken)
		api.PUT("/integrations/:id/settings", integrationsH.UpdateSettings)
		api.GET("/integrations/oauth/:provider", integrationsH.GetOAuthURL)
		api.POST("/integrations/oauth/:provider/callback", integrationsH.OAuthCallback)

		// Media
		api.GET("/media", mediaH.List)
		api.POST("/media/upload", mediaH.Upload)
		api.DELETE("/media/:id", mediaH.Delete)
		api.PATCH("/media/:id/alt", mediaH.UpdateAlt)
		api.POST("/media/profile-picture", mediaH.SetProfilePicture)

		// Analytics
		api.GET("/analytics", analyticsH.GetDashboard)
		api.GET("/analytics/channels", analyticsH.GetChannelOverview)
		api.GET("/analytics/schedule", analyticsH.GetScheduleAnalytics)
		api.GET("/analytics/integrations/:integrationId", analyticsH.GetIntegrationStats)
		api.GET("/analytics/posts/:postId", analyticsH.GetPostStats)

		// Tags
		api.GET("/tags", tagsH.List)
		api.POST("/tags", tagsH.Create)
		api.PUT("/tags/:id", tagsH.Update)
		api.DELETE("/tags/:id", tagsH.Delete)

		// Notifications
		api.GET("/notifications", notifH.List)
		api.POST("/notifications/read", notifH.MarkRead)
		api.GET("/notifications/unread", notifH.GetUnreadCount)

		// Webhooks
		api.GET("/webhooks", webhooksH.List)
		api.POST("/webhooks", webhooksH.Create)
		api.PUT("/webhooks/:id", webhooksH.Update)
		api.DELETE("/webhooks/:id", webhooksH.Delete)

		// Billing
		api.GET("/billing/subscription", billingH.GetSubscription)
		api.POST("/billing/checkout", billingH.CreateCheckout)
		api.POST("/billing/portal", billingH.CreatePortal)
		api.POST("/billing/cancel", billingH.CancelSubscription)
		api.GET("/billing/plans", billingH.GetPlans)

		// Settings
		api.GET("/settings", settingsH.GetOrgSettings)
		api.PUT("/settings", settingsH.UpdateOrgSettings)
		api.POST("/settings/api-key/regenerate", settingsH.RegenerateAPIKey)
		api.PATCH("/settings/shortlink", settingsH.UpdateShortLink)

		// Sets (templates)
		api.GET("/sets", setsH.List)
		api.POST("/sets", setsH.Create)
		api.PUT("/sets/:id", setsH.Update)
		api.DELETE("/sets/:id", setsH.Delete)

		// Signatures
		api.GET("/signatures", sigsH.List)
		api.POST("/signatures", sigsH.Create)
		api.PUT("/signatures/:id", sigsH.Update)
		api.DELETE("/signatures/:id", sigsH.Delete)

		// Auto-post
		api.GET("/auto-post", autoPostH.List)
		api.POST("/auto-post", autoPostH.Create)
		api.PUT("/auto-post/:id", autoPostH.Update)
		api.DELETE("/auto-post/:id", autoPostH.Delete)
		api.PATCH("/auto-post/:id/toggle", autoPostH.Toggle)

		// Announcements
		api.GET("/announcements", announcementsH.List)
		api.POST("/announcements", announcementsH.Create)
		api.DELETE("/announcements/:id", announcementsH.Delete)

		// OAuth App
		api.GET("/oauth-app", oauthAppH.Get)
		api.POST("/oauth-app", oauthAppH.Create)
		api.PUT("/oauth-app/:id", oauthAppH.Update)
		api.DELETE("/oauth-app/:id", oauthAppH.Delete)
	}

	// Public API v1 (API key auth)
	publicAPI := s.router.Group("/public-api/v1")
	publicAPI.Use(middleware.APIKeyAuth())
	{
		publicAPI.GET("/posts", postsH.List)
		publicAPI.POST("/posts", postsH.Create)
		publicAPI.GET("/integrations", integrationsH.List)
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
