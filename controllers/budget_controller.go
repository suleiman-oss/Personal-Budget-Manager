package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suleiman/Personal-Budget-Manager/models"
	"github.com/suleiman/Personal-Budget-Manager/services"
)

// CreateBudget creates a new budget
func CreateBudget(c *gin.Context) {
	var budget models.Budget
	budgetService := &services.BudgetService{}
	if err := c.BindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := budgetService.CreateBudget(&budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Budget created successfully", "budget": budget})
}

// GetBudgetByID retrieves a budget by ID
func GetBudgetByID(c *gin.Context) {
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	budgetService := &services.BudgetService{}
	budget, err := budgetService.GetBudgetByID(uint(budgetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// UpdateBudget updates an existing budget
func UpdateBudget(c *gin.Context) {
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	var budget models.Budget
	if err := c.BindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budgetService := &services.BudgetService{}
	if err := budgetService.UpdateBudget(uint(budgetID), &budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget updated successfully"})
}

// DeleteBudget deletes a budget by ID
func DeleteBudget(c *gin.Context) {
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	budgetService := &services.BudgetService{}
	if err := budgetService.DeleteBudget(uint(budgetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}
