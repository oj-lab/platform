package model

import "time"

type Problem struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Visibility  bool      `gorm:"not null"`
	CreateAt    time.Time `gorm:"autoCreateTime"`
	UpdateAt    time.Time `gorm:"autoUpdateTime"`
}

type TestCase struct {
	ID        string `gorm:"primaryKey"`
	ProblemID uint   `gorm:"not null"`
	Input     string `gorm:"not null"`
	Output    *string
}
