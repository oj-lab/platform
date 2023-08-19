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
