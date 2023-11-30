package database

import "time"

type Assignment struct {
	ID                 string       `gorm:"primaryKey;"`
	Name               string       `gorm:"type:varchar(100)" binding:"required" json:"name"`
	Points             int          `binding:"required,min=0,max=100" json:"points"`
	NumOfAttempts      int          `binding:"required,min=0,max=100" json:"num_of_attempts"`
	Deadline           time.Time    `binding:"required" gorm:"autoUpdateTime:false" json:"deadline"`
	Assignment_Created time.Time    `gorm:"autoUpdateTime:false"`
	Assignment_Updated time.Time    `gorm:"autoUpdateTime:false"`
	AccountID          string       `gorm:"foreignKey:AccountID" json:"-"`
	Submission         []Submission `json:"-"`
}
