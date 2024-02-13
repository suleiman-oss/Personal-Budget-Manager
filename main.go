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
	ucr := controllers.UserController{}
	ucr.UserService = &services.UserService{}
	ucr.UserService.DB = db
	// Set up routes
	setupRoutes(router, &ucr)

	// Run the server
	router.Run("localhost:8087")
}

func setupRoutes(router *gin.Engine, ucr *controllers.UserController) {

	router.POST("/register", ucr.CreateUser)
	router.POST("/login", ucr.Login)
	userRouter := router.Group("/users")
	{
		userRouter.PUT("/:id", ucr.UpdateUser)    // Update a user
		userRouter.DELETE("/:id", ucr.DeleteUser) // Delete a user
	}
}
