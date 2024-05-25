package user

import (
	"context"

	"github.com/google/uuid"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	"github.com/oj-lab/oj-lab-platform/modules/auth"
	"github.com/oj-lab/oj-lab-platform/modules/log"
)

func GetUser(ctx context.Context, account string) (*user_model.User, error) {
	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUser(db, account)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CheckUserExist(ctx context.Context, account string) (bool, error) {
	getOptions := user_model.GetUserOptions{
		Account: account,
	}
	db := gorm_agent.GetDefaultDB()
	count, err := user_model.CountUserByOptions(db, getOptions)
	if err != nil {
		return false, err
	}

	if count > 1 {
		log.AppLogger().Warnf("user %s has %d records", account, count)
	}

	return count > 0, nil
}

func StartLoginSession(ctx context.Context, account, password string) (*uuid.UUID, error) {
	db := gorm_agent.GetDefaultDB()
	user, err := user_model.GetUserByAccountPassword(db, account, password)
	if err != nil {
		return nil, err
	}

	loginSession := auth.NewLoginSession(auth.LoginSession{
		Account: account,
		Roles:   user.GetRolesStringArray(),
	})
	err = loginSession.SaveToRedis(ctx)
	if err != nil {
		return nil, err
	}

	return &loginSession.Id, nil
}
