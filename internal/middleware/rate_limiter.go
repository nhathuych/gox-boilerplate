package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

// NewRedisRateLimiter creates a Redis-backed rate limiting middleware using ulule/limiter.
// Keying strategy:
// - if authenticated user is present in context: key by user id
// - otherwise: key by client IP
type RateLimiterHandler gin.HandlerFunc

func NewRedisRateLimiter(rdb *redis.Client, limit int64, period time.Duration) RateLimiterHandler {
	// Ulule limiter creation can error; for a boilerplate we fall back to a no-op middleware.
	// (In production you'd return an error from fx.)
	store, err := redisstore.NewStore(rdb)
	if err != nil {
		return RateLimiterHandler(func(c *gin.Context) { c.Next() })
	}

	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}
	l := limiter.New(store, rate)

	return RateLimiterHandler(ginlimiter.NewMiddleware(
		l,
		ginlimiter.WithKeyGetter(func(c *gin.Context) string {
			if uid, ok := UserIDFromContext(c.Request.Context()); ok {
				return fmt.Sprintf("rl:user:%s", uid.String())
			}
			return fmt.Sprintf("rl:ip:%s", c.ClientIP())
		}),
	))
}
