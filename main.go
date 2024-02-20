// main.go

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
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
func GetDB() *gorm.DB {
	// Open connection to PostgreSQL
	driver, connectionString := ConnectionString()
	db, err := gorm.Open(driver, connectionString)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	//defer db.Close()

	// Check if the connection is successful
	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("Error pinging PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
	models.Migrate(db)
	return db
}
func main() {
	// Initialize Gin
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	db := GetDB()
	kafkaProducer := InitKafkaProducer()
	go runCronJob(db, kafkaProducer)
	cr := controllers.NewController(db)
	setupRoutes(router, cr)
	defer ClearSessionOnExit()
	router.Run("localhost:8087")
}
func ClearSessionOnExit() {
	// Perform an internal request to clear the session
	// You may use http.DefaultClient or any custom client you have.
	_, err := http.Get("http://localhost:8087/logout")
	if err != nil {
		fmt.Println("Failed to clear session on exit:", err)
	}
}
func runCronJob(db *gorm.DB, producer sarama.AsyncProducer) {
	c := cron.New()

	// Add a function to the cron scheduler using a closure to capture the db instance
	c.AddFunc("18 18 * * *", func() {
		services.CheckExpense(db, producer) // Pass the db instance to the checkExpenses function
	})

	// Start the cron scheduler
	c.Start()

	// Keep the main goroutine alive forever
	select {}
}
func InitKafkaProducer() sarama.AsyncProducer {
	producer, err := sarama.NewAsyncProducer([]string{"kafka-broker:9092"}, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return producer
}
func setupRoutes(router *gin.Engine, cr *controllers.Controller) {

	router.POST("/register", cr.UserController.CreateUser)
	router.POST("/login", cr.UserController.Login)
	router.GET("/logout", cr.UserController.Logout)
	userRouter := router.Group("/users")
	{
		userRouter.PUT("/update", cr.UserController.UpdateUser)    // Update a user
		userRouter.DELETE("/delete", cr.UserController.DeleteUser) // Delete a user
	}
	budgetRouter := router.Group("/budget")
	{
		budgetRouter.GET("/", cr.BudgetController.GetAllBudgets)
		budgetRouter.POST("/newbudget", cr.BudgetController.CreateBudget)
	}
	expenseRouter := router.Group("/expense")
	{
		expenseRouter.GET("/", cr.ExpenseController.GetAllExpenses)
		expenseRouter.POST("/newexpense", cr.ExpenseController.CreateExpense)
	}
}
