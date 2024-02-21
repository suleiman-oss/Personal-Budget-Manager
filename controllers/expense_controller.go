package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"github.com/suleiman-oss/Personal-Budget-Manager/services"
)

type ExpenseController struct {
	ExpenseService *services.ExpenseService
}

// CreateExpense creates a new expense
func (ecr *ExpenseController) CreateExpense(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	var expense models.Expense
	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	// Assign the user ID to the budget
	expense.UserID = userIDUint
	expense.Date = time.Now()
	if err := ecr.ExpenseService.CreateExpense(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "expense": expense})
}

// GetExpenseByID retrieves an expense by ID
func (ecr *ExpenseController) GetExpenseByID(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	expenseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	expense, err := ecr.ExpenseService.GetExpenseByID(uint(expenseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateExpense updates an existing expense
func (ecr *ExpenseController) UpdateExpense(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	expenseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	var expense models.Expense
	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ecr.ExpenseService.UpdateExpense(uint(expenseID), &expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully"})
}

// DeleteExpense deletes an expense by ID
func (ecr *ExpenseController) DeleteExpense(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	expenseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	if err := ecr.ExpenseService.DeleteExpense(uint(expenseID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

func (ecr *ExpenseController) GetAllExpenses(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	// Retrieve budgets for the authenticated user
	var expenses []models.Expense
	if err := ecr.ExpenseService.DB.Where("user_id = ?", userIDUint).Find(&expenses).Error; err != nil {
		// Handle error if any
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budgets"})
		return
	}

	// Return the list of budgets (even if it's an empty list)
	if len(expenses) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "User doesn't have any budget"})
		return
	}
	totalExpenses := 0.0
	for _, expense := range expenses {
		totalExpenses += expense.Amount
	}

	// Attach total expenses to the response
	response := gin.H{
		"expenses":      expenses,
		"totalExpenses": totalExpenses,
	}
	c.JSON(http.StatusOK, response)

}
