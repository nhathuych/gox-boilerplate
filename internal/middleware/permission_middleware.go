package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nhathuych/gox-boilerplate/internal/usecase"
)

// RequirePermission returns a Gin middleware that allows the request only if the
// authenticated user has the given permission (e.g. article:publish).
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, ok := PermissionsFromContext(c.Request.Context())
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing permissions"})
			return
		}
		if !usecase.HasPermission(perms, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
