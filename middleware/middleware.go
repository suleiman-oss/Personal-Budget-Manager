package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suleiman/Personal-Budget-Manager/models"
)

// AuthorizeRole checks if the user has the required role for accessing the resource
func AuthorizeRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Check if the user has the required role
		if user.(*models.User).Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
