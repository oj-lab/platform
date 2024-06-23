package user_model

import (
	"github.com/google/uuid"
	"github.com/oj-lab/oj-lab-platform/models"
)

type UserGroup struct {
	models.MetaFields
	UID          uuid.UUID         `json:"uid" gorm:"primaryKey"`
	OwnerAccount string            `json:"ownerAccount"`
	Owner        User              `json:"owner" gorm:"foreignKey:OwnerAccount"`
	Members      []UserGroupMember `json:"members" gorm:"foreignKey:GroupUID"`
}

type UserGroupMember struct {
	models.MetaFields
	GroupUID    uuid.UUID `json:"groupUID" gorm:"primaryKey"`
	UserAccount string    `json:"userAccount" gorm:"primaryKey"`
	User        User      `json:"user"`
	Role        string    `json:"role"`
}
