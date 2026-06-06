package middleware

import (
	"postpanda/backend-go/internal/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     []string{config.App.FrontendURL, "http://localhost:4200", "http://localhost:6274"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key", "x-org-id"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(cfg)
}
