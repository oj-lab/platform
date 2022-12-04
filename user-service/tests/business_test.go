package tests

import (
	"testing"

	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
)

func SetupTestDatabase() {
	dataBaseSettings, err := config.GetDatabaseSettings("../../config/ini/test.ini")
	if err != nil {
		panic("failed to get database settings")
	}
	business.OpenDBConnection(*dataBaseSettings)
}

func TestString2Role(t *testing.T) {
	role := model.String2Role("admin")
	if role != model.RoleAdmin {
		panic("wrong role")
	}
	role = model.String2Role("user")
	if role != model.RoleUser {
		panic("wrong role")
	}
	role = model.String2Role("wrong")
	if role != model.RoleUser {
		panic("wrong role")
	}
}

func TestCreateUser(t *testing.T) {
	SetupTestDatabase()
	_ = business.CreateUser("common", "password", []model.Role{model.RoleUser})
}

func TestDeleteUser(t *testing.T) {
	SetupTestDatabase()
	err := business.DeleteUser("common")
	if err != nil {
		panic(err)
	}
}

func TestComparePassword(t *testing.T) {
	SetupTestDatabase()
	_ = business.CreateUser("common", "password", []model.Role{model.RoleUser})
	_, err := business.ComparePassword("common", "password")
	if err != nil {
		panic(err)
	}
	match, _ := business.ComparePassword("common", "wrong password")
	if match {
		panic("wrong password should be rejected")
	}
}

func TestUpdateUser(t *testing.T) {
	SetupTestDatabase()
	_ = business.CreateUser("common", "password", []model.Role{model.RoleUser})
	newPassword := "new password"
	_ = business.UpdateUser("common", nil, &newPassword, nil, nil, nil)
	match, err := business.ComparePassword("common", "new password")
	if err != nil {
		panic(err)
	}
	if !match {
		panic("wrong password")
	}
}

func TestGetUserInfo(t *testing.T) {
	SetupTestDatabase()
	_ = business.CreateUser("common", "password", []model.Role{model.RoleUser})
	mobile := "12312312345"
	_ = business.UpdateUser("common", nil, nil, nil, nil, &mobile)
	userInfo, err := business.GetUserInfo(nil, nil, &mobile)
	if err != nil {
		panic(err)
	}
	if userInfo.Account != "common" {
		panic("wrong account")
	}
	account := "common"
	userInfo, err = business.GetUserInfo(&account, nil, nil)
	if err != nil {
		panic(err)
	}
	if userInfo.Account != "common" {
		panic("wrong account")
	}
}

func TestFindUserInfos(t *testing.T) {
	SetupTestDatabase()
	_ = business.CreateUser("common", "password", []model.Role{model.RoleUser})
	userInfo, err := business.FindUserInfos("%co%", 0, 1)
	if err != nil {
		panic(err)
	}
	if userInfo[0].Account != "common" {
		panic("wrong account")
	}
}

func TestCountUser(t *testing.T) {
	SetupTestDatabase()
	_ = business.CreateUser("common", "password", []model.Role{model.RoleUser})
	count, err := business.CountUser("%co%")
	if err != nil {
		panic(err)
	}
	if count == 0 {
		panic("wrong count")
	}
}
