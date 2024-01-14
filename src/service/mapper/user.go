package mapper

import (
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

// Account, Password, Roles will be used to create a new user.
func CreateUser(tx *gorm.DB, user model.User) error {
	hashedPassword, err := argon2id.CreateHash(*user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	User := model.User{
		Name:           user.Name,
		Account:        user.Account,
		HashedPassword: hashedPassword,
		Roles:          user.Roles,
	}

	return tx.Create(&User).Error
}

func GetUser(tx *gorm.DB, account string) (*model.User, error) {
	db_user := model.User{}
	err := tx.Model(&model.User{}).Preload("Roles").Where("account = ?", account).First(&db_user).Error
	if err != nil {
		return nil, err
	}

	return &db_user, err
}

func GetPublicUser(tx *gorm.DB, account string) (*model.User, error) {

	db_user := model.User{}
	err := tx.Model(&model.User{}).Preload("Roles").Select(model.PublicUserSelection).Where("account = ?", account).First(&db_user).Error
	if err != nil {
		return nil, err
	}

	return &db_user, err
}

func DeleteUser(tx *gorm.DB, user model.User) error {
	return tx.Delete(&model.User{Account: user.Account}).Error
}

func UpdateUser(tx *gorm.DB, update model.User) error {
	old := model.User{}
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
	if update.Roles != nil {
		new.Roles = update.Roles
	}

	return tx.Model(&model.User{Account: new.Account}).Updates(new).Error
}

type GetUserOptions struct {
	Account string
	Email   string
	Mobile  string
	Offset  *int
	Limit   *int
}

// Count the total number of users that match the options,
// ignoring the offset and limit.
func CountUserByOptions(tx *gorm.DB, options GetUserOptions) (int64, error) {
	var count int64

	tx = tx.
		Model(&model.User{}).
		Where("account = ?", options.Account).
		Or("email = ?", options.Email).
		Or("mobile = ?", options.Mobile)

	err := tx.Count(&count).Error

	return count, err
}

func GetUserByOptions(tx *gorm.DB, options GetUserOptions) ([]model.User, int64, error) {
	total, err := CountUserByOptions(tx, options)
	if err != nil {
		return nil, 0, err
	}

	db_users := []model.User{}

	tx = tx.
		Where("account = ?", options.Account).
		Or("email = ?", options.Email).
		Or("mobile = ?", options.Mobile)
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

func CheckUserPassword(tx *gorm.DB, account string, password string) (bool, error) {
	user := model.User{}
	err := tx.Where("account = ?", account).First(&user).Error
	if err != nil {
		return false, err
	}
	return argon2id.ComparePasswordAndHash(password, user.HashedPassword)
}
