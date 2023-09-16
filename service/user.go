package service

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/core"
	"github.com/OJ-lab/oj-lab-services/core/auth"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/sirupsen/logrus"
)

func GetUser(account string) (*model.User, *core.SeviceError) {
	user, err := mapper.GetUser(account)
	if err != nil {
		return nil, core.NewInternalError("failed to get user")
	}

	return user, nil
}

func CheckUserExist(account string) (bool, error) {
	getOptions := mapper.GetUserOptions{
		Account: account,
	}
	count, err := mapper.CountUserByOptions(getOptions)
	if err != nil {
		return false, err
	}

	if count > 1 {
		logrus.Warnf("user %s has %d records", account, count)
	}

	return count > 0, nil
}

func StartLoginSession(ctx context.Context, account, password string) (*string, *core.SeviceError) {
	match, err := mapper.CheckUserPassword(account, password)
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
