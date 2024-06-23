package user_model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUserGroup(tx *gorm.DB, group UserGroup) error {
	group.UID = uuid.New()

	return tx.Create(&group).Error
}

func CreateUserGroupMember(tx *gorm.DB, member UserGroupMember) error {
	return tx.Create(&member).Error
}
