package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant_id from context (set by auth middleware)
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tenant information not found"})
			c.Abort()
			return
		}

		// Set tenant_id in context for use in handlers
		c.Set("current_tenant_id", tenantID)
		c.Next()
	}
}
