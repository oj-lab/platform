package user_model

import "github.com/oj-lab/platform/models"

type User struct {
	models.MetaFields
	Account        string   `json:"account" gorm:"primaryKey"`
	Name           string   `json:"name"`
	Password       *string  `json:"password,omitempty" gorm:"-:all"`
	HashedPassword string   `json:"-" gorm:"not null"`
	Email          *string  `json:"email,omitempty" gorm:"unique"`
	AvatarURL      string   `json:"avatarUrl"`
	GithubLogin    *string  `json:"githubLogin" gorm:"unique"`
	Roles          []string `json:"roles" gorm:"-"`
}

func (u *User) IsAdmin() bool {
	for _, role := range u.Roles {
		if role == "admin" || role == "super" {
			return true
		}
	}
	return false
}

var PublicUserSelection = append([]string{"account", "name", "avatar_url"}, models.MetaFieldsSelection...)
