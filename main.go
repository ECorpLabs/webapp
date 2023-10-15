package main

import (
	"log"
	"os"
	nocache "webapp/api/handler"
	"webapp/controllers"
	database "webapp/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	// Connect to the database
	err := database.Connect()
	if err != nil {
		log.Println("Error connecting to the database")
	} else {
		database.Database.AutoMigrate(&database.Account{})
		database.Database.AutoMigrate(&database.Assignment{})

		// Seed the database
		filePath := os.Getenv("FILE_PATH")
		err = database.SeedData(database.Database, filePath)
		if err != nil {
			log.Println("Error seeding the database: ", err)
		}
	}
	// Create a router
	router := gin.Default()
	router.Use(nocache.NoCache())

	binding.EnableDecoderDisallowUnknownFields = true
	// Create a group for /healthz
	healthzGroup := router.Group("/healthz")
	{
		// Register health routes under /healthz
		controllers.RegisterHealthRoutes(healthzGroup)
	}
	// Create a group for authenticated users
	authGroup := router.Group("/v1/")
	{
		// Initialize AssignmentController and register its routes
		assignmentController := controllers.NewAssignmentController()
		assignmentController.RegisterRoutes(authGroup)
	}

	router.Run(":" + os.Getenv("APP_PORT"))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
