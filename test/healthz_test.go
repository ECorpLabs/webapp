package test

import (
	"log"
	"net/http/httptest"
	"testing"
	controller "webapp/controllers"
	"webapp/database"
	client "webapp/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HealthTestSuite struct {
	suite.Suite
	App *gin.Engine
}

func TestHealthTestSuite(t *testing.T) {
	suite.Run(t, &HealthTestSuite{})
}

func (s *HealthTestSuite) SetupSuite() {

	loadEnv()
	//get logger instance
	client.Init()
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	app := gin.New()
	//Connect to the database
	err := database.Connect()
	if err != nil {
		logger.Error("Error connecting to the database", zap.Error(err))
	}

	healthzGroup := app.Group("/healthz")
	{
		// Register health routes under /healthz
		controller.RegisterHealthRoutes(healthzGroup, logger)
	}
	s.App = app
}

func loadEnv() {
	err := godotenv.Load("../.env.local")
	if err != nil {
		log.Print("Error loading .env file", err)
	}
	log.Print("Environment variables loaded successfully!!")
}

func (s *HealthTestSuite) TestIntegrationHealth() {

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	s.App.ServeHTTP(w, req)

	res := w.Code

	s.Equal(200, res)
}
