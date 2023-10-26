package controllers

import (
	"net/http"
	"time"
	"webapp/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssignmentController struct {
	// Add any necessary services or dependencies here
}

func NewAssignmentController() *AssignmentController {
	// Initialize any services or dependencies here
	return &AssignmentController{}
}

func (ac *AssignmentController) RegisterRoutes(router *gin.RouterGroup) {
	unsupportedMethod := func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
	}
	// Group routes under "/assignments"
	assignmentsGroup := router.Group("/assignments")
	assignmentsGroup.PATCH("", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
	})
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
	account_id := c.GetString("account_id")
	var input database.Assignment

	// Bind the JSON data to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		//Returning assignment created response and return the assignment created
		c.JSON(http.StatusCreated, assignment)
		// c.JSON(http., gin.H{assignment : assignment})
	}
}

func (ac *AssignmentController) UpdateAssignment(c *gin.Context) {
	// Implement assignment update logic
	id := c.Param("id")
	account_id := c.GetString("account_id")

	// Check if assignment with given ID exists
	var existingAssignment database.Assignment
	if err := database.Database.First(&existingAssignment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if existingAssignment.AccountID != account_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var newAssignment database.Assignment
	if err := c.ShouldBindJSON(&newAssignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (ac *AssignmentController) DeleteAssignment(c *gin.Context) {
	// Implement assignment deletion logic
	id := c.Param("id")
	accountID := c.GetString("account_id") // Assuming account_id is set during authentication

	// Check if assignment with given ID exists
	var existingAssignment database.Assignment
	if err := database.Database.First(&existingAssignment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Check if the logged-in user is the creator of the assignment
	if existingAssignment.AccountID != accountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Delete the assignment
	if err := database.Database.Delete(&existingAssignment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ac *AssignmentController) GetAssignment(c *gin.Context) {
	// Implement assignment retrieval logic
	id := c.Param("id")

	// Check if assignment with given ID exists
	var assignment database.Assignment
	if err := database.Database.First(&assignment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func (ac *AssignmentController) GetAssignments(c *gin.Context) {
	// Implement fetching all assignments logic
	var assignments []database.Assignment
	if err := database.Database.Find(&assignments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, assignments)
	}
}
