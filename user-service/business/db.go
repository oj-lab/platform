package business

import (
	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/OJ-lab/oj-lab-services/packages/utils"
	"github.com/alexedwards/argon2id"
)

func CreateUser(account string, password string, roles model.Roles) error {
	db := application.GetDefaultDB()

	hashedPassword, err := utils.GetHashedPassword(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	user := model.UserTable{
		Account:        account,
		HashedPassword: hashedPassword,
		Roles:          roles.ToPQArray(),
	}

	return db.Create(&user).Error
}

func DeleteUser(account string) error {
	db := application.GetDefaultDB()
	return db.Delete(&model.UserTable{Account: account}).Error
}

func UpdateUser(account string, name *string, password *string, roles model.Roles, email *string, mobile *string) error {
	db := application.GetDefaultDB()
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

	user := model.UserTable{
		Account:        account,
		Name:           name,
		HashedPassword: hashedPassword,
		Roles:          roles.ToPQArray(),
		Email:          email,
		Mobile:         mobile,
	}

	return db.Model(&model.UserTable{Account: account}).Updates(user).Error
}

func ComparePassword(account string, password string) (bool, error) {
	db := application.GetDefaultDB()
	var user model.UserTable
	err := db.Where("account = ?", account).First(&user).Error
	if err != nil {
		return false, err
	}
	return argon2id.ComparePasswordAndHash(password, user.HashedPassword)
}

func GetUserInfo(maybeAccount *string, maybeEmail *string, maybeMobile *string) (*model.UserInfo, error) {
	db := application.GetDefaultDB()
	account := ""
	if maybeAccount != nil {
		account = *maybeAccount
	}
	var user model.UserTable
	err := db.Where(&model.UserTable{Account: account, Email: maybeEmail, Mobile: maybeMobile}).First(&user).Error
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

func FindUserInfos(query string, offset int, limit int) ([]model.UserInfo, error) {
	db := application.GetDefaultDB()
	var users []model.UserTable
	err := db.Where("account LIKE ?", query).Or("name LIKE ?", query).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	var userInfos []model.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, model.UserInfo{
			Account:  user.Account,
			Name:     user.Name,
			Roles:    model.PQArray2Roles(&user.Roles),
			Email:    user.Email,
			CreateAt: user.CreateAt,
			UpdateAt: user.UpdateAt,
		})
	}
	return userInfos, err
}

func CountUser(query string) (int64, error) {
	db := application.GetDefaultDB()
	var count int64
	err := db.Model(&model.UserTable{}).Where("account LIKE ?", query).Or("name LIKE ?", query).Count(&count).Error
	return count, err
}
