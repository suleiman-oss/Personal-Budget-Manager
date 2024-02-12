// main.go

package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/suleiman/Personal-Budget-Manager/config"
	"github.com/suleiman/Personal-Budget-Manager/controllers"
)

// ConnectionString returns the connection string for PostgreSQL
func ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
}

func main() {
	// Open connection to PostgreSQL
	db, err := sql.Open("postgres", ConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
}

func setupRoutes(router *gin.Engine) {
	// Routes for user management
	userRouter := router.Group("/users")
	{
		userRouter.POST("/", controllers.CreateUser)      // Create a new user
		userRouter.GET("/:id", controllers.GetUserByID)   // Get a user by ID
		userRouter.PUT("/:id", controllers.UpdateUser)    // Update a user
		userRouter.DELETE("/:id", controllers.DeleteUser) // Delete a user
	}

	// Routes for expense management
	expenseRouter := router.Group("/expenses")
	{
		expenseRouter.POST("/", controllers.CreateExpense)      // Create a new expense
		expenseRouter.GET("/:id", controllers.GetExpenseByID)   // Get an expense by ID
		expenseRouter.PUT("/:id", controllers.UpdateExpense)    // Update an expense
		expenseRouter.DELETE("/:id", controllers.DeleteExpense) // Delete an expense
	}

	// Routes for budget management
	budgetRouter := router.Group("/budgets")
	{
		budgetRouter.POST("/", controllers.CreateBudget)      // Create a new budget
		budgetRouter.GET("/:id", controllers.GetBudgetByID)   // Get a budget by ID
		budgetRouter.PUT("/:id", controllers.UpdateBudget)    // Update a budget
		budgetRouter.DELETE("/:id", controllers.DeleteBudget) // Delete a budget
	}

	// Routes for notification management
	notificationRouter := router.Group("/notifications")
	{
		notificationRouter.POST("/", controllers.CreateNotification)      // Create a new notification
		notificationRouter.GET("/:id", controllers.GetNotificationByID)   // Get a notification by ID
		notificationRouter.PUT("/:id", controllers.UpdateNotification)    // Update a notification
		notificationRouter.DELETE("/:id", controllers.DeleteNotification) // Delete a notification
	}
}
