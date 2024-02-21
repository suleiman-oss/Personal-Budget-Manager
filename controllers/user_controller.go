package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"github.com/suleiman-oss/Personal-Budget-Manager/services"
)

type UserController struct {
	UserService *services.UserService
}

// CreateUser creates a new user
func (ucr *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ucr.UserService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
	//c.Redirect(http.StatusAccepted, "/login")
}

// GetUserByID retrieves a user by ID
func (ucr *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	i, _ := strconv.Atoi(userID)
	user, err := ucr.UserService.GetUserByID(uint(i))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates an existing user
func (ucr *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	i, _ := strconv.Atoi(userID)
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ucr.UserService.UpdateUser(uint(i), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user by ID
func (ucr *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	i, _ := strconv.Atoi(userID)
	if err := ucr.UserService.DeleteUser(uint(i)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ucr *UserController) Login(c *gin.Context) {
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
	user, err := ucr.UserService.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		c.Abort()
		return
	}

	// Verify password
	if !ucr.UserService.VerifyPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		c.Abort()
		return
	}
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
func (ucr *UserController) Logout(c *gin.Context) {
	// Get session
	session := sessions.Default(c)

	// Clear session data
	session.Clear()

	// Save the session to apply changes
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Session cleared successfully"})
}
