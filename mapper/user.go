package mapper

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/database"
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/alexedwards/argon2id"
)

func CreateUser(ctx context.Context, user model.User) error {
	db := database.GetDefaultDB()
	hashedPassword, err := utils.GetHashedPassword(*user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	userTable := model.UserTable{
		Account:        user.Account,
		HashedPassword: hashedPassword,
		Roles:          user.Roles.ToPQArray(),
	}

	return db.Create(&userTable).Error
}

func DeleteUser(ctx context.Context, user model.User) error {
	db := database.GetDefaultDB()
	return db.Delete(&model.UserTable{Account: user.Account}).Error
}

func UpdateUser(ctx context.Context, user model.User) error {
	db := database.GetDefaultDB()
	var hashedPassword string
	if user.Password != nil {
		var err error
		hashedPassword, err = utils.GetHashedPassword(*user.Password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
	} else {
		hashedPassword = ""
	}

	userRow := model.UserTable{
		Account:        user.Account,
		Name:           user.Name,
		HashedPassword: hashedPassword,
		Roles:          user.Roles.ToPQArray(),
		Email:          user.Email,
		Mobile:         user.Mobile,
	}

	return db.Model(&model.UserTable{Account: userRow.Account}).Updates(userRow).Error
}

type UserOption struct {
	Account *string
	Email   *string
	Mobile  *string
}

func GetUserByOption(option UserOption) (*model.UserInfo, error) {
	db := database.GetDefaultDB()
	account := ""
	if option.Account != nil {
		account = *option.Account
	}
	var user model.UserTable
	err := db.Where(&model.UserTable{Account: account, Email: option.Email, Mobile: option.Mobile}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &model.UserInfo{
		Account:  user.Account,
		Name:     user.Name,
		Roles:    model.PQArray2Roles(&user.Roles),
		Email:    user.Email,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}, nil
}
