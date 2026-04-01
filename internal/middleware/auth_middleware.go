package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nhathuych/gox-boilerplate/internal/auth"
)

func Auth(jwt *auth.JWTService, blacklist *auth.TokenBlacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimPrefix(strings.TrimSpace(h), "Bearer ")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		claims, err := jwt.ParseAccess(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		jti := claims.ID
		bl, err := blacklist.IsBlacklisted(ctx, jti)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth check failed"})
			return
		}
		if bl {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		c.Request = c.Request.WithContext(WithUser(ctx, claims.UserID, claims.Permissions, jti))
		c.Next()
	}
}
