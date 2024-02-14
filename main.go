// main.go

package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/suleiman/Personal-Budget-Manager/config"
	"github.com/suleiman/Personal-Budget-Manager/controllers"
	"github.com/suleiman/Personal-Budget-Manager/models"
	"github.com/suleiman/Personal-Budget-Manager/services"
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
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
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
	cr := controllers.Controller{}
	ucr := controllers.UserController{}
	ucr.UserService = &services.UserService{}
	ucr.UserService.DB = db
	bcr := controllers.BudgetController{}
	bcr.BudgetService = &services.BudgetService{}
	bcr.BudgetService.DB = db
	cr.UserController = &ucr
	cr.BudgetController = &bcr
	// Set up routes
	setupRoutes(router, &cr)

	// Run the server
	router.Run("localhost:8087")
}

func setupRoutes(router *gin.Engine, cr *controllers.Controller) {

	router.POST("/register", cr.UserController.CreateUser)
	router.POST("/login", cr.UserController.Login)
	userRouter := router.Group("/users")
	{
		userRouter.PUT("/:id", cr.UserController.UpdateUser)    // Update a user
		userRouter.DELETE("/:id", cr.UserController.DeleteUser) // Delete a user
	}
	budgetRouter := router.Group("/budget")
	{
		budgetRouter.GET("/", cr.BudgetController.GetAllBudgets)
		budgetRouter.POST("/newbudget", cr.BudgetController.CreateBudget)
	}
}
