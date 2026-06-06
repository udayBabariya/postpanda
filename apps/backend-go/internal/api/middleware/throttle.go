package middleware

import (
	"context"
	"fmt"
	"net/http"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/database"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("throttle:%s", ip)
		ttl := time.Duration(config.App.ThrottleTTL) * time.Second
		limit := config.App.ThrottleLimit

		ctx := context.Background()
		count, err := database.Redis.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			database.Redis.Expire(ctx, key, ttl)
		}

		if int(count) > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message":    "Too Many Requests",
				"statusCode": http.StatusTooManyRequests,
			})
			return
		}

		c.Next()
	}
}
