package user

import "github.com/oj-lab/oj-lab-platform/models"

type User struct {
	models.MetaFields
	Account        string  `gorm:"primaryKey" json:"account"`
	Name           string  `json:"name"`
	Password       *string `gorm:"-:all" json:"password,omitempty"`
	HashedPassword string  `gorm:"not null" json:"-"`
	Roles          []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Email          *string `gorm:"unique" json:"email,omitempty"`
	Mobile         *string `gorm:"unique" json:"mobile,omitempty"`
}

var PublicUserSelection = append([]string{"account", "name"}, models.MetaFieldsSelection...)

type Role struct {
	models.MetaFields
	Name  string  `gorm:"primaryKey" json:"name"`
	Users []*User `gorm:"many2many:user_roles" json:"users,omitempty"`
}

func (user User) GetRolesStringSet() map[string]struct{} {
	roleSet := map[string]struct{}{}
	for _, role := range user.Roles {
		roleSet[role.Name] = struct{}{}
	}
	return roleSet
}
