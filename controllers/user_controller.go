package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suleiman/Personal-Budget-Manager/models"
	"github.com/suleiman/Personal-Budget-Manager/services"
)

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user models.User
	userService := &services.UserService{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

// GetUserByID retrieves a user by ID
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	i, _ := strconv.Atoi(userID)
	userService := &services.UserService{}
	user, err := userService.GetUserByID(uint(i))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates an existing user
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	i, _ := strconv.Atoi(userID)
	var user models.User
	userService := &services.UserService{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userService.UpdateUser(uint(i), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	i, _ := strconv.Atoi(userID)
	userService := &services.UserService{}
	if err := userService.DeleteUser(uint(i)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
