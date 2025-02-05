package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimitMiddleware returns a Gin middleware function that limits requests.
func RateLimitMiddleware(rdb *redis.Client, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// Increment the request counter.
		count, err := rdb.Incr(c, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			return
		}

		// Set expiration on the key if it's the first request.
		if count == 1 {
			rdb.Expire(c, key, time.Minute)
		}

		// If the count exceeds the limit, abort with 429 status.
		if count > int64(requestsPerMinute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": "60 seconds",
			})
			return
		}

		// Optionally, set a header with remaining requests.
		c.Header("X-RateLimit-Remaining", strconv.Itoa(requestsPerMinute-int(count)))
		c.Next()
	}
}
