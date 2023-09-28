package mapper

import (
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/alexedwards/argon2id"
)

// Account, Password, Roles will be used to create a new user.
func CreateUser(user model.User) error {
	db := gormAgent.GetDefaultDB()
	hashedPassword, err := argon2id.CreateHash(*user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	User := model.User{
		Account:        user.Account,
		HashedPassword: hashedPassword,
		Roles:          user.Roles,
	}

	return db.Create(&User).Error
}

func GetUser(account string) (*model.User, error) {
	db := gormAgent.GetDefaultDB()
	db_user := model.User{}
	err := db.Model(&model.User{}).Preload("Roles").Where("account = ?", account).First(&db_user).Error
	if err != nil {
		return nil, err
	}

	return &db_user, err
}

func GetPublicUser(account string) (*model.User, error) {
	db := gormAgent.GetDefaultDB()
	db_user := model.User{}
	err := db.Model(&model.User{}).Preload("Roles").Select(model.PublicUserSelection).Where("account = ?", account).First(&db_user).Error
	if err != nil {
		return nil, err
	}

	return &db_user, err
}

func DeleteUser(user model.User) error {
	db := gormAgent.GetDefaultDB()
	return db.Delete(&model.User{Account: user.Account}).Error
}

func UpdateUser(update model.User) error {
	db := gormAgent.GetDefaultDB()

	old := model.User{}
	err := db.Where("account = ?", update.Account).First(&old).Error
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

	return db.Model(&model.User{Account: new.Account}).Updates(new).Error
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
func CountUserByOptions(options GetUserOptions) (int64, error) {
	db := gormAgent.GetDefaultDB()
	var count int64

	tx := db.
		Model(&model.User{}).
		Where("account = ?", options.Account).
		Or("email = ?", options.Email).
		Or("mobile = ?", options.Mobile)

	err := tx.Count(&count).Error

	return count, err
}

func GetUserByOptions(options GetUserOptions) ([]model.User, int64, error) {
	total, err := CountUserByOptions(options)
	if err != nil {
		return nil, 0, err
	}

	db := gormAgent.GetDefaultDB()
	db_users := []model.User{}

	tx := db.
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

func CheckUserPassword(account string, password string) (bool, error) {
	db := gormAgent.GetDefaultDB()
	user := model.User{}
	err := db.Where("account = ?", account).First(&user).Error
	if err != nil {
		return false, err
	}
	return argon2id.ComparePasswordAndHash(password, user.HashedPassword)
}
