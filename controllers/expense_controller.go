package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suleiman/Personal-Budget-Manager/models"
	"github.com/suleiman/Personal-Budget-Manager/services"
)

// CreateExpense creates a new expense
func CreateExpense(c *gin.Context) {
	var expense models.Expense
	expenseService := &services.ExpenseService{}
	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := expenseService.CreateExpense(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "expense": expense})
}

// GetExpenseByID retrieves an expense by ID
func GetExpenseByID(c *gin.Context) {
	expenseID := c.Param("id")
	expenseService := &services.ExpenseService{}
	i, _ := strconv.Atoi(expenseID)
	expense, err := expenseService.GetExpenseByID(uint(i))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateExpense updates an existing expense
func UpdateExpense(c *gin.Context) {
	expenseID := c.Param("id")
	i, _ := strconv.Atoi(expenseID)
	var expense models.Expense
	expenseService := &services.ExpenseService{}
	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := expenseService.UpdateExpense(uint(i), &expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully"})
}

// DeleteExpense deletes an expense by ID
func DeleteExpense(c *gin.Context) {
	expenseID := c.Param("id")
	i, _ := strconv.Atoi(expenseID)
	expenseService := &services.ExpenseService{}
	if err := expenseService.DeleteExpense(uint(i)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
