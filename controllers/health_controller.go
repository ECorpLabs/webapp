package controllers

import (
	"net/http"

	"webapp/database"
	client "webapp/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterHealthRoutes(group *gin.RouterGroup, logger *zap.Logger) {

	logRequest := func(c *gin.Context) {
		logger.Info("Request received",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
		)
		c.Next()
		logger.Info("Response sent",
			zap.Int("status", c.Writer.Status()),
		)
	}

	logError := func(c *gin.Context, statusCode int, errorMessage string) {
		logger.Error(errorMessage,
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
		)
		c.JSON(statusCode, gin.H{"error": errorMessage})
	}

	unsupportedMethod := func(c *gin.Context) {
		logError(c, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	group.GET("", logRequest, func(c *gin.Context) {

		if c.Request.Body != http.NoBody || len(c.Request.URL.Query()) > 0 {
			logError(c, http.StatusBadRequest, "Status Bad Request")
			return
		}
		// err := database.Connect()
		// ping the database to check connection
		err := database.Database.Exec("SELECT 1").Error
		if err != nil {
			logError(c, http.StatusServiceUnavailable, "Status Service Unavailable")
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		client.GetMetricsClient().Incr("web.get", 1)
	})

	group.POST("", logRequest, func(c *gin.Context) {
		client.GetMetricsClient().Incr("web.post", 1)
		unsupportedMethod(c)
	})
	group.PUT("", logRequest, func(c *gin.Context) {
		client.GetMetricsClient().Incr("web.put", 1)
		unsupportedMethod(c)
	})
	group.DELETE("", logRequest, func(c *gin.Context) {
		client.GetMetricsClient().Incr("web.delete", 1)
		unsupportedMethod(c)
	})
	group.PATCH("", logRequest, func(c *gin.Context) {
		client.GetMetricsClient().Incr("web.patch", 1)
		unsupportedMethod(c)
	})
}
