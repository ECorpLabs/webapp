package controllers

import (
	"net/http"

	"webapp/database"

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
		err := database.Connect()
		if err != nil {
			logError(c, http.StatusServiceUnavailable, "Status Service Unavailable")
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
	})

	group.POST("", logRequest, unsupportedMethod)
	group.PUT("", logRequest, unsupportedMethod)
	group.DELETE("", logRequest, unsupportedMethod)
	group.PATCH("", logRequest, unsupportedMethod)
}
