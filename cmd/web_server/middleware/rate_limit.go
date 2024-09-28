package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	redis_agent "github.com/oj-lab/platform/modules/agent/redis"
)

const rateLimitKeyFormat = "RL_%s_%s" // "RL_<route>_<ip>"

func BuildHandleRateLimitWithDuration(duration time.Duration) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		redisClient := redis_agent.GetDefaultRedisClient()
		key := fmt.Sprintf(rateLimitKeyFormat, ginCtx.FullPath(), ginCtx.ClientIP())
		result := redisClient.SetNX(ginCtx, key, 1, duration)
		if !result.Val() {
			ginCtx.String(http.StatusTooManyRequests, "Too frequent requests, please try again later")
			ginCtx.Abort()
			return
		}
		ginCtx.Next()
	}
}
