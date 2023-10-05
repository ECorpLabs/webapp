package database

import "time"

type Account struct {
	ID             string       `gorm:"primaryKey;"`
	FirstName      string       `binding:"required" json:"first_name"`
	LastName       string       `binding:"required" json:"last_name"`
	Email          string       `binding:"required" gorm:"unique_index"`
	Password       string       `binding:"required" json:"-"` // Hide password from JSON response
	AccountCreated time.Time    `gorm:"autoUpdateTime:false"`
	AccountUpdated time.Time    `gorm:"autoUpdateTime:false"`
	Assignments    []Assignment `json:"-"`
}
