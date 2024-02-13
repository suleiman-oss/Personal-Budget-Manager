package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suleiman/Personal-Budget-Manager/models"
	"github.com/suleiman/Personal-Budget-Manager/services"
)

// AuthMiddleware handles user authentication and authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := models.LoginInput{}
		// Bind and validate input
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if email and password are provided
		if input.Email == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
			c.Abort()
			return
		}

		// Authenticate user
		userService := &services.UserService{}
		user, err := userService.GetUserByEmail(input.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
			c.Abort()
			return
		}

		// Verify password
		if !services.VerifyPassword(user.Password, input.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			c.Abort()
			return
		}

		// Attach user information to the context
		c.Set("user", user)

		// Proceed to the next handler
		c.Next()
	}
}

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
