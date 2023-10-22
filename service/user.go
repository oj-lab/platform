package service

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/core"
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/core/auth"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func GetUser(ctx context.Context, account string) (*model.User, *core.SeviceError) {
	db := gormAgent.GetDefaultDB()
	user, err := mapper.GetUser(db, account)
	if err != nil {
		return nil, core.NewInternalError("failed to get user")
	}

	return user, nil
}

func CheckUserExist(ctx context.Context, account string) (bool, error) {
	getOptions := mapper.GetUserOptions{
		Account: account,
	}
	db := gormAgent.GetDefaultDB()
	count, err := mapper.CountUserByOptions(db, getOptions)
	if err != nil {
		return false, err
	}

	if count > 1 {
		core.AppLogger().Warnf("user %s has %d records", account, count)
	}

	return count > 0, nil
}

func StartLoginSession(ctx context.Context, account, password string) (*string, *core.SeviceError) {
	db := gormAgent.GetDefaultDB()
	match, err := mapper.CheckUserPassword(db, account, password)
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}
	if !match {
		return nil, core.NewUnauthorizedError("invalid account or password")
	}

	loginSession := auth.NewLoginSession(account)
	err = loginSession.SaveToRedis(ctx)
	if err != nil {
		return nil, core.NewInternalError(err.Error())
	}

	return &loginSession.Id, nil
}
