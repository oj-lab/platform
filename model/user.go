package model

import (
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/alexedwards/argon2id"
	"time"
)

type User struct {
	Account        string `gorm:"primaryKey"`
	Name           *string
	HashedPassword string    `gorm:"not null"`
	Role           string    `gorm:"not null"`
	Email          *string   `gorm:"unique"`
	Mobile         *string   `gorm:"unique"`
	CreateAt       time.Time `gorm:"autoCreateTime"`
	UpdateAt       time.Time `gorm:"autoUpdateTime"`
}

type UserInfo struct {
	Account  string
	Name     *string
	Role     string
	Email    *string
	Mobile   *string
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

func CreateUser(account string, password string, role Role) error {
	hashedPassword, err := utils.GetHashedPassword(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	user := User{
		Account:        account,
		HashedPassword: hashedPassword,
		Role:           string(role),
	}

	return db.Create(&user).Error
}

func DeleteUser(account string) error {
	return db.Delete(&User{Account: account}).Error
}

func UpdateUser(account string, name *string, password *string, role *Role, email *string, mobile *string) error {
	var hashedPassword string
	if password != nil {
		var err error
		hashedPassword, err = utils.GetHashedPassword(*password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
	} else {
		hashedPassword = ""
	}

	roleString := ""
	if role != nil {
		roleString = string(*role)
	}

	user := User{
		Account:        account,
		Name:           name,
		HashedPassword: hashedPassword,
		Role:           roleString,
		Email:          email,
		Mobile:         mobile,
	}

	return db.Model(&User{Account: account}).Updates(user).Error
}

func ComparePassword(account string, password string) (bool, error) {
	var user User
	err := db.Where("account = ?", account).First(&user).Error
	if err != nil {
		return false, err
	}
	return argon2id.ComparePasswordAndHash(password, user.HashedPassword)
}

func GetUserInfo(maybeAccount *string, maybeEmail *string, maybeMobile *string) (*UserInfo, error) {
	account := ""
	if maybeAccount != nil {
		account = *maybeAccount
	}
	var user User
	err := db.Where(&User{Account: account, Email: maybeEmail, Mobile: maybeMobile}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &UserInfo{
		Account:  user.Account,
		Name:     user.Name,
		Role:     user.Role,
		Email:    user.Email,
		Mobile:   user.Mobile,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}, nil
}

func FindUserInfos(query string, offset int, limit int) ([]UserInfo, error) {
	var users []User
	err := db.Where("account LIKE ?", query).Or("name LIKE ?", query).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	var userInfos []UserInfo
	for _, user := range users {
		userInfos = append(userInfos, UserInfo{
			Account:  user.Account,
			Name:     user.Name,
			Role:     user.Role,
			Email:    user.Email,
			Mobile:   user.Mobile,
			CreateAt: user.CreateAt,
			UpdateAt: user.UpdateAt,
		})
	}
	return userInfos, err
}

func CountUser(query string) (int64, error) {
	var count int64
	err := db.Model(&User{}).Where("account LIKE ?", query).Or("name LIKE ?", query).Count(&count).Error
	return count, err
}
