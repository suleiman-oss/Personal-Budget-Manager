package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"github.com/suleiman-oss/Personal-Budget-Manager/services"
)

type GoalController struct {
	GoalService *services.GoalService
}

// NewGoalController creates a new instance of GoalController
func NewGoalController(goalService *services.GoalService) *GoalController {
	return &GoalController{GoalService: goalService}
}

// CreateGoalHandler handles the creation of a new goal
func (gc *GoalController) CreateGoalHandler(c *gin.Context) {
	var goal models.Goal
	if err := c.BindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := gc.GoalService.CreateGoal(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": goal})
}

// GetGoalHandler handles the retrieval of a goal by ID
func (gc *GoalController) GetGoalHandler(c *gin.Context) {
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}
	goal, err := gc.GoalService.GetGoalByID(uint(goalID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": *goal})
}

// UpdateGoalHandler handles the updating of an existing goal
func (gc *GoalController) UpdateGoalHandler(c *gin.Context) {
	var updatedGoal models.Goal
	if err := c.BindJSON(&updatedGoal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := gc.GoalService.UpdateGoal(&updatedGoal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedGoal})
}

// DeleteGoalHandler handles the deletion of an existing goal
func (gc *GoalController) DeleteGoalHandler(c *gin.Context) {
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}
	if err := gc.GoalService.DeleteGoal(uint(goalID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}
