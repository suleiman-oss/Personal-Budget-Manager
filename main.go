// main.go

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"github.com/suleiman-oss/Personal-Budget-Manager/config"
	"github.com/suleiman-oss/Personal-Budget-Manager/controllers"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"github.com/suleiman-oss/Personal-Budget-Manager/services"
	"gopkg.in/Shopify/sarama.v1"
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
	router.Use(CORSMiddleware())
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	db := GetDB()
	kafkaProducer := InitKafkaProducer()
	go runCronJob(db, kafkaProducer)
	cr := controllers.NewController(db)
	setupRoutes(router, cr)
	defer ClearSessionOnExit()
	//router.Run("localhost:8087")
	config, err := config.LoadConfig("utils/")
	if err != nil {
		log.Fatal("cannot load config:", err.Error())
	}
	port := config.Port
	if port == "" {
		port = "localhost:8087"
	}
	if err := router.Run(port); err != nil {
		log.Panicf("error: %s", err)
	}
}
func ClearSessionOnExit() {
	// Perform an internal request to clear the session
	// You may use http.DefaultClient or any custom client you have.
	_, err := http.Get("http://localhost:8087/logout")
	if err != nil {
		fmt.Println("Failed to clear session on exit:", err)
	}
}
func runCronJob(db *gorm.DB, producer sarama.SyncProducer) {
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
func InitKafkaProducer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
