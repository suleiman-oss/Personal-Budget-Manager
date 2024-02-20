package services

import (
	"time"

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
func CheckExpense(db *gorm.DB) {
	s := NotificationService{}
	s.DB = db
	now := time.Now()
	currentMonth := now.Format("2006-01")
	var users []models.User
	s.DB.Find(&users)
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		var expenses []models.Expense
		var budgets []models.Budget
		var totalExpenses, totalBudget float64

		s.DB.Where("user_id = ? AND date >= ? AND date < ?", user.ID, currentMonth, now.AddDate(0, 1, 0)).Find(&expenses)
		s.DB.Where("user_id = ? AND month = ?", user.ID, string(currentMonth)).Find(&budgets)

		for _, expense := range expenses {
			totalExpenses += expense.Amount
		}

		for _, budget := range budgets {
			totalBudget += budget.Limit
		}

		if totalExpenses == user.Income+user.PreviousSavings-totalBudget {
			notification := models.Notification{
				Message: "Your expenses for the month have reached your budget limit.",
				Type:    "Budget Limit Reached",
				UserID:  user.ID,
			}
			s.CreateNotification(&notification)
		}
	}
}
