package user_model

import "github.com/oj-lab/oj-lab-platform/models"

type User struct {
	models.MetaFields
	Account        string  `json:"account" gorm:"primaryKey"`
	Name           string  `json:"name"`
	Password       *string `json:"password,omitempty" gorm:"-:all"`
	HashedPassword string  `json:"-" gorm:"not null"`
	Roles          []*Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	Email          *string `json:"email,omitempty" gorm:"unique"`
	Mobile         *string `json:"mobile,omitempty" gorm:"unique"`
}

var PublicUserSelection = append([]string{"account", "name"}, models.MetaFieldsSelection...)

type Role struct {
	models.MetaFields
	Name  string  `json:"name" gorm:"primaryKey"`
	Users []*User `json:"users,omitempty" gorm:"many2many:user_roles"`
}

func (user User) GetRolesStringSet() map[string]struct{} {
	roleSet := map[string]struct{}{}
	for _, role := range user.Roles {
		roleSet[role.Name] = struct{}{}
	}
	return roleSet
}
