package service

import (
	"github.com/OJ-lab/oj-lab-services/package/mapper"
	"github.com/sirupsen/logrus"
)

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
