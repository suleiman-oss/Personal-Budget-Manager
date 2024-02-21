package services

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
)

// BudgetService provides CRUD operations for the Budget model
type BudgetService struct {
	DB *gorm.DB
}

// CreateBudget creates a new budget
func (s *BudgetService) CreateBudget(budget *models.Budget) error {
	return s.DB.Create(budget).Error
}

// GetBudgetByID retrieves a budget by ID
func (s *BudgetService) GetBudgetByID(budgetID uint) (*models.Budget, error) {
	var budget models.Budget
	err := s.DB.First(&budget, budgetID).Error
	return &budget, err
}

// UpdateBudget updates an existing budget
func (s *BudgetService) UpdateBudget(budgetID uint, budget *models.Budget) error {
	existingBudget, err := s.GetBudgetByID(budgetID)
	if err != nil {
		return err
	}
	existingBudget.Category = budget.Category
	existingBudget.Limit = budget.Limit
	existingBudget.UserID = budget.UserID

	return s.DB.Save(budget).Error
}

// DeleteBudget deletes a budget by ID
func (s *BudgetService) DeleteBudget(budgetID uint) error {
	return s.DB.Delete(&models.Budget{}, budgetID).Error
}
