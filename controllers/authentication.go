// authentication.go

package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"webapp/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// DecodeBase64Hash decodes the base64-encoded hash to get username and password
func DecodeBase64Hash(encodedHash string) (username, password string, err error) {
	// Remove the "Basic " prefix from the hash
	encodedHash = strings.TrimPrefix(encodedHash, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return "", "", err
	}
	decodedHash := string(decodedBytes)
	creds := strings.Split(decodedHash, ":")
	if len(creds) != 2 {
		return "", "", errors.New("invalid base64-encoded hash")
	}

	return creds[0], creds[1], nil
}

// AuthMiddleware implements authentication logic
func AuthMiddleware(c *gin.Context) {
	// Get the base64-encoded hash from the request header
	encodedHash := c.GetHeader("Authorization")
	// println(encodedHash)

	// Decode the base64-encoded hash to get username and password
	username, password, err := DecodeBase64Hash(encodedHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64-encoded hash"})
		c.Abort()
		return
	}

	// Retrieve user from the database based on username
	var account database.Account
	if err := database.Database.Where("email = ?", username).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username"})
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		c.Abort()
		return
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
	c.Set("account_id", account.ID)

	// Authentication successful
	c.Next()
}
