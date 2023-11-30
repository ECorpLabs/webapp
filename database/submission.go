package database

import "time"

type Submission struct {
	ID                 string    `gorm:"primaryKey;"`
	Assignment_Id      string    `json:"assignment_id"`
	SubmissionUrl      string    `binding:"required" json:"submission_url"`
	Submission_Date    time.Time `gorm:"autoUpdateTime:false"`
	Submission_Updated time.Time `gorm:"autoUpdateTime:false"`
	AccountID          string    `gorm:"foreignKey:AccountID" json:"-"`
}
