package model

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	Account        string `gorm:"primaryKey"`
	Name           *string
	HashedPassword string         `gorm:"not null"`
	Roles          pq.StringArray `gorm:"not null;type:varchar(255)[]"`
	Email          *string        `gorm:"unique"`
	Mobile         *string        `gorm:"unique"`
	CreateAt       time.Time      `gorm:"autoCreateTime"`
	UpdateAt       time.Time      `gorm:"autoUpdateTime"`
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

type Roles []Role

func (roles *Roles) ToPQArray() pq.StringArray {
	res := pq.StringArray{}
	for _, role := range *roles {
		res = append(res, string(role))
	}
	return res
}

func PQArray2Roles(roleArray *pq.StringArray) Roles {
	res := Roles{}
	for _, rolePQString := range *roleArray {
		res = append(res, String2Role(rolePQString))
	}
	return res
}

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
