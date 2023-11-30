package main

import (
	"log"
	"os"
	nocache "webapp/api/handler"
	"webapp/controllers"
	database "webapp/database"
	client "webapp/logger"

	zap "go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	// Initialize the metrics logger
	client.Init()

	// Create a logger
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		os.Getenv("LOG_FILE_PATH"),
	}
	config.DisableStacktrace = true
	logger := zap.Must(config.Build())
	defer logger.Sync()

	// logger := zap.Must(zap.NewProduction())
	// defer logger.Sync()

	// Connect to the database
	err := database.Connect()
	if err != nil {
		logger.Error("Error connecting to the database", zap.Error(err))
	} else {
		database.Database.AutoMigrate(&database.Account{})
		database.Database.AutoMigrate(&database.Assignment{})
		database.Database.AutoMigrate(&database.Submission{})

		// Seed the database
		filePath := os.Getenv("FILE_PATH")
		err = database.SeedData(database.Database, filePath)
		if err != nil {
			logger.Error("Error seeding the database", zap.Error(err))
		}
	}
	// Create a router
	router := gin.Default()
	router.Use(nocache.NoCache())

	router.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		// metricsClient.Incr("web.request", 1)
		c.Next()
	})

	binding.EnableDecoderDisallowUnknownFields = true
	// Create a group for /healthz
	healthzGroup := router.Group("/healthz")
	{
		// Register health routes under /healthz

		controllers.RegisterHealthRoutes(healthzGroup, logger)
	}
	// Create a group for authenticated users

	authGroup := router.Group("/v1/")
	{
		// Initialize AssignmentController and register its routes
		assignmentController := controllers.NewAssignmentController()
		assignmentController.RegisterRoutes(authGroup, logger)
	}

	router.Run(":" + os.Getenv("APP_PORT"))
	defer client.GetMetricsClient().Close()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
