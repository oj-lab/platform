package user_model

import (
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/oj-lab/oj-lab-platform/models"
	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Account, Password, Roles will be used to create a new user.
func CreateUser(tx *gorm.DB, request User) (*User, error) {
	user := User{
		MetaFields: models.NewMetaFields(),
		Name:       request.Name,
		Account:    request.Account,
		Email:      request.Email,
		AvatarURL:  request.AvatarURL,
	}

	if request.Password != nil {
		hashedPassword, err := argon2id.CreateHash(*request.Password, argon2id.DefaultParams)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = hashedPassword
	}

	err := tx.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(tx *gorm.DB, account string) (*User, error) {
	db_user := User{}
	err := tx.Model(&User{}).Where("account = ?", account).First(&db_user).Error
	if err != nil {
		return nil, err
	}

	return &db_user, err
}

func GetPublicUser(tx *gorm.DB, account string) (*User, error) {
	db_user := User{}
	err := tx.Model(&User{}).Select(PublicUserSelection).
		Where("account = ?", account).First(&db_user).Error
	if err != nil {
		return nil, err
	}

	return &db_user, err
}

func DeleteUser(tx *gorm.DB, user User) error {
	return tx.Select(clause.Associations).Delete(&User{Account: user.Account}).Error
}

func UpdateUser(tx *gorm.DB, update User) error {
	old := User{}
	err := tx.Where("account = ?", update.Account).First(&old).Error
	if err != nil {
		return err
	}

	hashedPassword := ""
	if update.Password != nil {
		hashedPassword, err = argon2id.CreateHash(*update.Password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
	}

	new := old
	if update.Password != nil {
		new.HashedPassword = hashedPassword
	}

	return tx.Model(&User{Account: new.Account}).Updates(new).Error
}

type GetUserOptions struct {
	AccountQuery string
	EmailQuery   string
	DomainRole   *casbin_agent.DomainRole
	Offset       *int
	Limit        *int
}

// Count the total number of users that match the options,
// ignoring the offset and limit.
func CountUserByOptions(tx *gorm.DB, options GetUserOptions) (int64, error) {
	var count int64

	tx = tx.Model(&User{})

	if options.AccountQuery != "" {
		tx = tx.Where("account LIKE ?", options.AccountQuery)
	}
	if options.EmailQuery != "" {
		tx = tx.Where("email LIKE ?", options.EmailQuery)
	}

	err := tx.Count(&count).Error

	return count, err
}

func GetUserByOptions(tx *gorm.DB, options GetUserOptions) ([]User, int64, error) {
	total, err := CountUserByOptions(tx, options)
	if err != nil {
		return nil, 0, err
	}

	db_users := []User{}

	if options.AccountQuery != "" {
		tx = tx.Where("account LIKE ?", options.AccountQuery)
	}
	if options.EmailQuery != "" {
		tx = tx.Where("email LIKE ?", options.EmailQuery)
	}

	if options.Offset != nil {
		tx = tx.Offset(*options.Offset)
	}
	if options.Limit != nil {
		tx = tx.Limit(*options.Limit)
	}

	err = tx.Find(&db_users).Error
	if err != nil {
		return nil, 0, err
	}

	return db_users, total, nil
}

func GetUserByAccountPassword(tx *gorm.DB, account string, password string) (*User, error) {
	user := User{}
	err := tx.Where("account = ?", account).First(&user).Error
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("password not match")
	}

	return &user, nil
}
