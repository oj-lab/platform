package business

import (
	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/model"
	user_utils "github.com/OJ-lab/oj-lab-services/user-service/utils"
	"github.com/OJ-lab/oj-lab-services/utils"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDBConnection(settings config.DatabaseSettings) {
	db = utils.MustGetDBConnection(settings)
}

func CreateUser(account string, password string, roles model.Roles) error {
	hashedPassword, err := user_utils.GetHashedPassword(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	user := model.User{
		Account:        account,
		HashedPassword: hashedPassword,
		Roles:          roles.ToPQArray(),
	}

	return db.Create(&user).Error
}

func DeleteUser(account string) error {
	return db.Delete(&model.User{Account: account}).Error
}

func UpdateUser(account string, name *string, password *string, roles model.Roles, email *string, mobile *string) error {
	var hashedPassword string
	if password != nil {
		var err error
		hashedPassword, err = user_utils.GetHashedPassword(*password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
	} else {
		hashedPassword = ""
	}

	user := model.User{
		Account:        account,
		Name:           name,
		HashedPassword: hashedPassword,
		Roles:          roles.ToPQArray(),
		Email:          email,
		Mobile:         mobile,
	}

	return db.Model(&model.User{Account: account}).Updates(user).Error
}

func ComparePassword(account string, password string) (bool, error) {
	var user model.User
	err := db.Where("account = ?", account).First(&user).Error
	if err != nil {
		return false, err
	}
	return argon2id.ComparePasswordAndHash(password, user.HashedPassword)
}

func GetUserInfo(maybeAccount *string, maybeEmail *string, maybeMobile *string) (*model.UserInfo, error) {
	account := ""
	if maybeAccount != nil {
		account = *maybeAccount
	}
	var user model.User
	err := db.Where(&model.User{Account: account, Email: maybeEmail, Mobile: maybeMobile}).First(&user).Error
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
	var users []model.User
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
	var count int64
	err := db.Model(&model.User{}).Where("account LIKE ?", query).Or("name LIKE ?", query).Count(&count).Error
	return count, err
}
