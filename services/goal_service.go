package services

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
)

type GoalService struct {
	DB *gorm.DB
}

// NewGoalService creates a new instance of GoalService
func NewGoalService(db *gorm.DB) *GoalService {
	return &GoalService{DB: db}
}

// CreateGoal creates a new goal
func (gs *GoalService) CreateGoal(goal *models.Goal) error {
	if err := gs.DB.Create(goal).Error; err != nil {
		return err
	}
	return nil
}

// GetGoalByID retrieves a goal by its ID
func (gs *GoalService) GetGoalByID(id uint) (*models.Goal, error) {
	var goal models.Goal
	if err := gs.DB.First(&goal, id).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

// UpdateGoal updates an existing goal
func (gs *GoalService) UpdateGoal(goal *models.Goal) error {
	if err := gs.DB.Save(*goal).Error; err != nil {
		return err
	}
	return nil
}

// DeleteGoal deletes a goal
func (gs *GoalService) DeleteGoal(id uint) error {
	if err := gs.DB.Where("id = ?", id).Delete(&models.Goal{}).Error; err != nil {
		return err
	}
	return nil
}

// GetGoalsByUserID retrieves all goals for a given user ID
func (gs *GoalService) GetGoalsByUserID(userID uint) ([]models.Goal, error) {
	var goals []models.Goal
	if err := gs.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}
