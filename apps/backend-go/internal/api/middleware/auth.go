package middleware

import (
	"net/http"
	"postpanda/backend-go/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey   = "userId"
	OrgIDKey    = "organizationId"
	ClaimsKey   = "claims"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(OrgIDKey, claims.OrganizationID)
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "API key required"})
			return
		}
		c.Set("apiKey", apiKey)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearer := c.GetHeader("Authorization")
	if bearer != "" && strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}

	cookie, err := c.Cookie("postiz-auth")
	if err == nil && cookie != "" {
		return cookie
	}

	return ""
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

func GetOrgID(c *gin.Context) string {
	return c.GetString(OrgIDKey)
}
