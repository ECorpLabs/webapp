package controllers

import (
	"net/http"

	"webapp/database"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(group *gin.RouterGroup) {
	unsupportedMethod := func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
	}

	group.GET("", func(c *gin.Context) {
		if c.Request.Body != http.NoBody || len(c.Request.URL.Query()) > 0 {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		err := database.Connect()
		if err != nil {
			c.Writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
	})

	group.POST("", unsupportedMethod)
	group.PUT("", unsupportedMethod)
	group.DELETE("", unsupportedMethod)
	group.PATCH("", unsupportedMethod)
}
