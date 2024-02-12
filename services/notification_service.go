package services

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman/Personal-Budget-Manager/models"
)

// NotificationService provides CRUD operations for the Notification model
type NotificationService struct {
	DB *gorm.DB
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	return s.DB.Create(notification).Error
}

// GetNotificationByID retrieves a notification by ID
func (s *NotificationService) GetNotificationByID(notificationID uint) (*models.Notification, error) {
	var notification models.Notification
	err := s.DB.First(&notification, notificationID).Error
	return &notification, err
}

// UpdateNotification updates an existing notification
func (s *NotificationService) UpdateNotification(notification *models.Notification) error {
	return s.DB.Save(notification).Error
}

// DeleteNotification deletes a notification by ID
func (s *NotificationService) DeleteNotification(notificationID uint) error {
	return s.DB.Delete(&models.Notification{}, notificationID).Error
}
