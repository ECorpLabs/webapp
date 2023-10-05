package database

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB, filePath string) error {
	// Read the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip the header line assuming it's present
	for _, line := range lines[1:] {
		// Parse data from CSV
		firstName := line[0]
		lastName := line[1]
		email := line[2]
		password := line[3]

		// Check if user already exists
		var existingUser Account
		result := db.Where("email = ?", email).First(&existingUser)
		if result.Error == nil {
			// User already exists, skip adding this user
			continue
		}

		// Hash the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		id := uuid.New().String()

		currentTime := time.Now().UTC()

		user := Account{
			ID:             id,
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			Password:       string(hashedPassword),
			AccountCreated: currentTime,
			AccountUpdated: currentTime,
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
