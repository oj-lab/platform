package model

import (
	"time"
)

type User struct {
	Account        string `gorm:"primaryKey"`
	Name           *string
	HashedPassword string    `gorm:"not null"`
	Roles          []Role    `gorm:"not null"`
	Email          *string   `gorm:"unique"`
	Mobile         *string   `gorm:"unique"`
	CreateAt       time.Time `gorm:"autoCreateTime"`
	UpdateAt       time.Time `gorm:"autoUpdateTime"`
}

type UserInfo struct {
	Account  string
	Name     *string
	Roles    []Role
	Email    *string
	CreateAt time.Time
	UpdateAt time.Time
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func String2Role(s string) Role {
	switch s {
	case "admin":
		return RoleAdmin
	case "user":
		return RoleUser
	default:
		return RoleUser
	}
}
