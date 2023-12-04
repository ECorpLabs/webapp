package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"webapp/database"
	"webapp/sns"

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

func GetUserEmailByID(accountID string) (string, error) {
	var account database.Account
	if err := database.Database.Where("id = ?", accountID).First(&account).Error; err != nil {
		return "", errors.New("user not found")
	}
	return account.Email, nil
}

var status string

func isZIP(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	status = "INVALIDZIP"
	contentType := resp.Header.Get("Content-Type")
	return strings.EqualFold(contentType, "application/zip")
	//check content length is not zero
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
	assignmentsGroup.POST("/:id/submission", ac.SubmitAssignment)
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
		ID:                 id,
		Name:               input.Name,
		Points:             input.Points,
		NumOfAttempts:      input.NumOfAttempts,
		Deadline:           input.Deadline,
		Assignment_Created: time.Now().UTC(),
		Assignment_Updated: time.Now().UTC(),
		AccountID:          account_id,
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
	existingAssignment.Assignment_Updated = time.Now().UTC()

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

	zap.L().Info("Assignment retrieved successfully", zap.String("id", id))
	c.JSON(http.StatusOK, existingAssignment)
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

func (ac *AssignmentController) SubmitAssignment(c *gin.Context) {
	// Implement assignment submission logic
	client.GetMetricsClient().Incr("web.post", 1)
	assignment_id := c.Param("id")
	account_id := c.GetString("account_id")
	var input database.Submission
	var assignment database.Assignment

	userEmail, err := GetUserEmailByID(account_id)
	if err != nil {
		logError(c, http.StatusInternalServerError, "Failed to fetch user email")
		return
	}

	var status string
	id := uuid.New().String()

	defer func() {
		//publish message to sns with the user details
		userDetails := struct {
			SubmissionID  string `json:"submission_id"`
			UserEmail     string `json:"user_email"`
			NumOfAttempts int    `json:"number_of_attempts"`
			SubmissionUrl string `json:"submission_url"`
			AssignmentID  string `json:"assignment_id"`
			Status        string `json:"status"`
		}{
			SubmissionID:  id, // Assuming `id` is the submission ID
			UserEmail:     userEmail,
			NumOfAttempts: assignment.NumOfAttempts, // User's email fetched from the database
			SubmissionUrl: input.SubmissionUrl,
			AssignmentID:  assignment_id,
			Status:        status,
		}
		// Convert userDetails to JSON
		userDetailsJSON, err := json.Marshal(userDetails)
		if err != nil {
			logError(c, http.StatusInternalServerError, "Failed to prepare user details")
			return
		}

		// Publish the message to SNS
		snsErr := sns.PublishToSNS(string(userDetailsJSON))
		if snsErr != nil {
			logError(c, http.StatusInternalServerError, snsErr.Error())
		}
	}()
	// Bind the JSON data to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		status = "INVALIDJSON"
		return
	}
	if !isZIP(input.SubmissionUrl) {
		logError(c, http.StatusBadRequest, "Invalid file format or empty ZIP file. Please submit a ZIP file url with content.")
		status = "INVALIDZIP"
		return
	}

	var count int64
	if err := database.Database.Model(&database.Submission{}).
		Where("assignment_id = ? AND account_id = ?", assignment_id, account_id).
		Count(&count).Error; err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	}

	// var assignment database.Assignment
	if err := database.Database.First(&assignment, "id = ?", assignment_id).Error; err != nil {
		logError(c, http.StatusNotFound, "Assignment not found")
		status = "ASSNOTFOUND"
		return
	}

	if count >= int64(assignment.NumOfAttempts) {
		logError(c, http.StatusForbidden, "Number of attempts exceeded for this assignment")
		status = "EXCEEDEDATTEMPTS"
		return
	}

	deadlineTime, _ := time.Parse(time.RFC3339, assignment.Deadline.Format(time.RFC3339))

	if time.Now().UTC().After(deadlineTime) {
		logError(c, http.StatusForbidden, "Assignment submission deadline has passed")
		status = "DEADLINEPASSED"
		return
	}

	submission := database.Submission{
		ID:                 id,
		Assignment_Id:      assignment_id,
		SubmissionUrl:      input.SubmissionUrl,
		Submission_Date:    time.Now().UTC(),
		Submission_Updated: time.Now().UTC(),
		AccountID:          account_id,
	}

	if err := database.Database.Save(&submission).Error; err != nil {
		logError(c, http.StatusBadRequest, err.Error())
		return
	} else {
		//Returning assignment created response and return the assignment created
		zap.L().Info("Submission Accepted", zap.String("id", id))
		c.JSON(http.StatusCreated, submission)
		status = "SUBMISSIONACCEPTED"
	}

	// //publish message to sns with the user details
	// userDetails := struct {
	// 	SubmissionID  string `json:"submission_id"`
	// 	UserEmail     string `json:"user_email"`
	// 	NumOfAttempts int    `json:"number_of_attempts"`
	// 	SubmissionUrl string `json:"submission_url"`
	// 	AssignmentID  string `json:"assignment_id"`
	// 	Status        string `json:"status"`
	// }{
	// 	SubmissionID:  id, // Assuming `id` is the submission ID
	// 	UserEmail:     userEmail,
	// 	NumOfAttempts: assignment.NumOfAttempts, // User's email fetched from the database
	// 	SubmissionUrl: input.SubmissionUrl,
	// 	AssignmentID:  assignment_id,
	// 	Status:        status,
	// }
	// // Convert userDetails to JSON
	// userDetailsJSON, err := json.Marshal(userDetails)
	// if err != nil {
	// 	logError(c, http.StatusInternalServerError, "Failed to prepare user details")
	// 	return
	// }

	// defer func() {
	// 	// Publish the message to SNS
	// 	snsErr := sns.PublishToSNS(string(userDetailsJSON))
	// 	if snsErr != nil {
	// 		logError(c, http.StatusInternalServerError, snsErr.Error())
	// 	}
	// }()

}
