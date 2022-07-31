package model

import (
	"time"
)

type User struct {
	Account        string `gorm:"primaryKey"`
	Name           string
	HashedPassword string    `gorm:"not null"`
	Role           string    `gorm:"not null"`
	Email          string    `gorm:"unique"`
	Mobile         string    `gorm:"unique"`
	CreateAt       time.Time `gorm:"autoCreateTime"`
	UpdateAt       time.Time `gorm:"autoUpdateTime"`
	DeleteAt       time.Time
}
