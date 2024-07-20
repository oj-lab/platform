package user_service

import (
	"context"

	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	"github.com/oj-lab/oj-lab-platform/modules/auth"
	log_module "github.com/oj-lab/oj-lab-platform/modules/log"
)

func GetUser(ctx context.Context, account string) (*user_model.User, error) {
	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUser(db, account)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserList(
	ctx context.Context, options user_model.GetUserOptions,
) ([]user_model.User, int64, error) {
	db := gorm_agent.GetDefaultDB()
	users, total, err := user_model.GetUserByOptions(db, options)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func UpdateUser(ctx context.Context, user user_model.User) error {
	db := gorm_agent.GetDefaultDB()
	err := user_model.UpdateUser(db, user)
	if err != nil {
		return err
	}

	return auth.UpdateLoginSessionByAccount(ctx,
		user.Account,
		auth.LoginSessionData{
			RoleSet: user.GetRolesStringSet(),
		})
}

func CheckUserExist(ctx context.Context, account string) (bool, error) {
	getOptions := user_model.GetUserOptions{
		AccountQuery: account,
	}
	db := gorm_agent.GetDefaultDB()
	count, err := user_model.CountUserByOptions(db, getOptions)
	if err != nil {
		return false, err
	}

	if count > 1 {
		log_module.AppLogger().
			WithField("account", account).
			WithField("count", count).
			Warn("user account is not unique")
	}

	return count > 0, nil
}

func StartLoginSession(ctx context.Context, account, password string) (*auth.LoginSession, error) {
	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUserByAccountPassword(db, account, password)
	if err != nil {
		return nil, err
	}

	ls := auth.NewLoginSession(account, auth.LoginSessionData{
		RoleSet: user.GetRolesStringSet(),
	})
	err = ls.SaveToRedis(ctx)
	if err != nil {
		return nil, err
	}

	return ls, nil
}
