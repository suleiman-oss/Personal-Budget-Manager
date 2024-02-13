// main.go

package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/suleiman/Personal-Budget-Manager/config"
	"github.com/suleiman/Personal-Budget-Manager/controllers"
	"github.com/suleiman/Personal-Budget-Manager/middleware"
	"github.com/suleiman/Personal-Budget-Manager/models"
)

// ConnectionString returns the connection string for PostgreSQL
func ConnectionString() (string, string) {
	config, err := config.LoadConfig("utils/")
	if err != nil {
		log.Fatal("cannot load config:", err.Error())
	}
	return config.DBDriver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
}

func main() {
	// Initialize Gin
	router := gin.Default()
	// Open connection to PostgreSQL
	driver, connectionString := ConnectionString()
	db, err := gorm.Open(driver, connectionString)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("Error pinging PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
	models.Migrate(db)

	// Set up routes
	setupRoutes(router)

	// Run the server
	router.Run("localhost:8087")
}

func setupRoutes(router *gin.Engine) {

	router.Use(middleware.AuthMiddleware())
	// Routes for user management
	// auth := router.Group("/auth")
	// {
	// 	// Register and login routes
	// 	auth.POST("/register", controllers.Register)
	// 	auth.POST("/login", controllers.login)
	// }
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
