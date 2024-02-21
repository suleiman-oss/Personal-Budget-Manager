package services

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
)

// ExpenseService provides CRUD operations for the Expense model
type ExpenseService struct {
	DB *gorm.DB
}

// CreateExpense creates a new expense
func (s *ExpenseService) CreateExpense(expense *models.Expense) error {
	return s.DB.Create(expense).Error
}

// GetExpenseByID retrieves an expense by ID
func (s *ExpenseService) GetExpenseByID(expenseID uint) (*models.Expense, error) {
	var expense models.Expense
	err := s.DB.First(&expense, expenseID).Error
	return &expense, err
}

// UpdateExpense updates an existing expense
func (s *ExpenseService) UpdateExpense(expenseID uint, expense *models.Expense) error {
	existingExpense, err := s.GetExpenseByID(expenseID)
	if err != nil {
		return err
	}

	existingExpense.Description = expense.Description
	existingExpense.Amount = expense.Amount
	existingExpense.Date = expense.Date
	existingExpense.Category = expense.Category
	existingExpense.UserID = expense.UserID

	return s.DB.Save(expense).Error
}

// DeleteExpense deletes an expense by ID
func (s *ExpenseService) DeleteExpense(expenseID uint) error {
	return s.DB.Delete(&models.Expense{}, expenseID).Error
}
