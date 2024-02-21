package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"github.com/suleiman-oss/Personal-Budget-Manager/services"
)

// CreateNotification creates a new notification
func CreateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notificationService := &services.NotificationService{}
	if err := notificationService.CreateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully", "notification": notification})
}

// GetNotificationByID retrieves a notification by ID
func GetNotificationByID(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	notificationService := &services.NotificationService{}
	notification, err := notificationService.GetNotificationByID(uint(notificationID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// UpdateNotification updates an existing notification
func UpdateNotification(c *gin.Context) {
	// notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
	// 	return
	// }

	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notificationService := &services.NotificationService{}
	if err := notificationService.UpdateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification updated successfully"})
}

// DeleteNotification deletes a notification by ID
func DeleteNotification(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	notificationService := &services.NotificationService{}
	if err := notificationService.DeleteNotification(uint(notificationID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
