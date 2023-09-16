package model

import "time"

type Submission struct {
	ID         string    `gorm:"primaryKey"`
	ProblemID  string    `gorm:"not null"`
	UserID     string    `gorm:"not null"`
	Language   string    `gorm:"not null"`
	Code       string    `gorm:"not null"`
	Status     string    `gorm:"not null"`
	Visibility bool      `gorm:"not null"`
	CreateAt   time.Time `gorm:"autoCreateTime"`
	UpdateAt   time.Time `gorm:"autoUpdateTime"`
}

type CheckPoint struct {
	ID           string `gorm:"primaryKey"`
	SubmissionID string `gorm:"not null"`
	TestCaseID   string `gorm:"not null"`
	Result       string `gorm:"not null"`
}
