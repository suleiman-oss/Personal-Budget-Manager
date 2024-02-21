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

type BudgetController struct {
	BudgetService *services.BudgetService
}

// CreateBudget creates a new budget
func (bcr *BudgetController) CreateBudget(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}

	// Bind JSON payload to Budget struct
	var budget models.Budget
	if err := c.BindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	// Assign the user ID to the budget
	budget.UserID = userIDUint
	now := time.Now()
	currentMonth := now.Format("2006-01")
	budget.Month = string(currentMonth)

	// Create the budget in the database
	if err := bcr.BudgetService.DB.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Budget created successfully", "budget": budget})
}

// GetBudgetByID retrieves a budget by ID
func (bcr *BudgetController) GetBudgetByID(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	budget, err := bcr.BudgetService.GetBudgetByID(uint(budgetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (bcr *BudgetController) GetAllBudgets(c *gin.Context) {
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
	var budgets []models.Budget
	if err := bcr.BudgetService.DB.Where("user_id = ?", userIDUint).Find(&budgets).Error; err != nil {
		// Handle error if any
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budgets"})
		return
	}

	// Return the list of budgets (even if it's an empty list)
	if len(budgets) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "User doesn't have any budget"})
	} else {
		c.JSON(http.StatusOK, budgets)
	}
}

// UpdateBudget updates an existing budget
func (bcr *BudgetController) UpdateBudget(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
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
	if err := bcr.BudgetService.UpdateBudget(uint(budgetID), &budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget updated successfully"})
}

// DeleteBudget deletes a budget by ID
func (bcr *BudgetController) DeleteBudget(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// User is not authenticated, return unauthorized status
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not logged in"})
		return
	}
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}
	if err := bcr.BudgetService.DeleteBudget(uint(budgetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}
