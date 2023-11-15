package controllers

import (
	"net/http"
	"time"
	"webapp/database"

	client "webapp/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AssignmentController struct {
	// Add any necessary services or dependencies here
}

func NewAssignmentController() *AssignmentController {
	// Initialize any services or dependencies here
	return &AssignmentController{}
}
func logError(c *gin.Context, statusCode int, errorMessage string) {
	logger := c.MustGet("logger").(*zap.Logger)
	logger.Error(errorMessage,
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.Int("status", c.Writer.Status()),
	)
	c.JSON(statusCode, gin.H{"error": errorMessage})

}

func (ac *AssignmentController) RegisterRoutes(router *gin.RouterGroup, logger *zap.Logger) {

	logRequest := func(c *gin.Context) {
		logger.Info("Request received",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Int("status", c.Writer.Status()),
		)
		c.Next()
		logger.Info("Response sent",
			zap.Int("status", c.Writer.Status()),
		)
	}
	// logError := func(c *gin.Context, statusCode int, errorMessage string) {
	// 	logger.Error(errorMessage,
	// 		zap.String("method", c.Request.Method),
	// 		zap.String("url", c.Request.URL.String()),
	// 	)
	// 	c.JSON(statusCode, gin.H{"error": errorMessage})
	// }

	unsupportedMethod := func(c *gin.Context) {
		logError(c, http.StatusMethodNotAllowed, "Method Not Allowed")
		client.GetMetricsClient().Incr("web.patch", 1)
	}
	// Group routes under "/assignments"
	assignmentsGroup := router.Group("/assignments")
	assignmentsGroup.Use(logRequest)
	assignmentsGroup.PATCH("", unsupportedMethod)
	// Add middleware for authentication
	assignmentsGroup.Use(AuthMiddleware)

	// Define assignment routes
	assignmentsGroup.POST("", ac.CreateAssignment)
	assignmentsGroup.PATCH("/:id", unsupportedMethod)
	assignmentsGroup.PUT("/:id", ac.UpdateAssignment)
	assignmentsGroup.DELETE("/:id", ac.DeleteAssignment)
	assignmentsGroup.GET("/:id", ac.GetAssignment)
	assignmentsGroup.GET("", ac.GetAssignments)

}

func (ac *AssignmentController) CreateAssignment(c *gin.Context) {
	// Implement assignment creation logic
	// Example: c.JSON(http.StatusOK, gin.H{"message": "Assignment created successfully"})
	client.GetMetricsClient().Incr("web.post", 1)
	account_id := c.GetString("account_id")
	var input database.Assignment

	// Bind the JSON data to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	}
	id := uuid.New().String()
	assignment := database.Assignment{
		ID:                id,
		Name:              input.Name,
		Points:            input.Points,
		NumOfAttempts:     input.NumOfAttempts,
		Deadline:          input.Deadline,
		AssignmentCreated: time.Now().UTC(),
		AssignmentUpdated: time.Now().UTC(),
		AccountID:         account_id,
	}
	if err := database.Database.Save(&assignment).Error; err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	} else {
		//Returning assignment created response and return the assignment created
		zap.L().Info("Assignment created successfully", zap.String("id", id))
		c.JSON(http.StatusCreated, assignment)

	}
}

func (ac *AssignmentController) UpdateAssignment(c *gin.Context) {
	// Implement assignment update logic
	client.GetMetricsClient().Incr("web.put", 1)
	id := c.Param("id")
	account_id := c.GetString("account_id")

	// Check if assignment with given ID exists
	var existingAssignment database.Assignment
	if err := database.Database.First(&existingAssignment, "id = ?", id).Error; err != nil {
		logError(c, http.StatusNotFound, "Assignment not found")
		return
	}

	if existingAssignment.AccountID != account_id {
		logError(c, http.StatusForbidden, "Permission denied: User does not own this assignment")
		return
	}

	var newAssignment database.Assignment
	if err := c.ShouldBindJSON(&newAssignment); err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update assignment fields
	existingAssignment.Name = newAssignment.Name
	existingAssignment.Points = newAssignment.Points
	existingAssignment.NumOfAttempts = newAssignment.NumOfAttempts
	existingAssignment.Deadline = newAssignment.Deadline
	existingAssignment.AssignmentUpdated = time.Now().UTC()

	// Save the updated assignment
	if err := database.Database.Save(&existingAssignment).Error; err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	}
	zap.L().Info("Assignment updated successfully", zap.String("id", id))
	c.Status(http.StatusNoContent)
}

func (ac *AssignmentController) DeleteAssignment(c *gin.Context) {
	// Implement assignment deletion logic
	client.GetMetricsClient().Incr("web.delete", 1)
	id := c.Param("id")
	accountID := c.GetString("account_id") // Assuming account_id is set during authentication

	// Check if assignment with given ID exists
	var existingAssignment database.Assignment
	if err := database.Database.First(&existingAssignment, "id = ?", id).Error; err != nil {
		logError(c, http.StatusNotFound, "Assignment not found")
		return
	}

	// Check if the logged-in user is the creator of the assignment
	if existingAssignment.AccountID != accountID {
		logError(c, http.StatusForbidden, "Permission denied: User does not own this assignment")
		return
	}

	// Delete the assignment
	if err := database.Database.Delete(&existingAssignment).Error; err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	}
	zap.L().Info("Assignment deleted successfully", zap.String("id", id))
	c.Status(http.StatusNoContent)
}

func (ac *AssignmentController) GetAssignment(c *gin.Context) {
	// Implement assignment retrieval logic
	client.GetMetricsClient().Incr("web.get", 1)
	id := c.Param("id")

	// Check if assignment with given ID exists
	var assignment database.Assignment
	if err := database.Database.First(&assignment, "id = ?", id).Error; err != nil {
		logError(c, http.StatusNotFound, "Assignment not found")
		return
	}
	zap.L().Info("Assignment retrieved successfully", zap.String("id", id))
	c.JSON(http.StatusOK, assignment)
}

func (ac *AssignmentController) GetAssignments(c *gin.Context) {
	// Implement fetching all assignments logic
	client.GetMetricsClient().Incr("web.get", 1)
	var assignments []database.Assignment
	if err := database.Database.Find(&assignments).Error; err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	} else {
		zap.L().Info("Assignments retrieved successfully")
		c.JSON(http.StatusOK, assignments)
	}
}
