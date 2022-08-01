package user_tests

import (
	"github.com/OJ-lab/oj-lab-services/model"
	"testing"
)

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
	_ = model.CreateUser("common", "password", model.RoleUser)
}

func TestDeleteUser(t *testing.T) {
	SetupTestDatabase()
	err := model.DeleteUser("common")
	if err != nil {
		panic(err)
	}
}

func TestComparePassword(t *testing.T) {
	SetupTestDatabase()
	_ = model.CreateUser("common", "password", model.RoleUser)
	match, err := model.ComparePassword("common", "password")
	if err != nil {
		panic(err)
	}
	match, err = model.ComparePassword("common", "wrong password")
	if match {
		panic("wrong password should be rejected")
	}
}

func TestUpdateUser(t *testing.T) {
	SetupTestDatabase()
	_ = model.CreateUser("common", "password", model.RoleUser)
	newPassword := "new password"
	_ = model.UpdateUser("common", nil, &newPassword, nil, nil, nil)
	match, err := model.ComparePassword("common", "new password")
	if err != nil {
		panic(err)
	}
	if !match {
		panic("wrong password")
	}
}

func TestGetUserInfo(t *testing.T) {
	SetupTestDatabase()
	_ = model.CreateUser("common", "password", model.RoleUser)
	mobile := "12312312345"
	_ = model.UpdateUser("common", nil, nil, nil, nil, &mobile)
	userInfo, err := model.GetUserInfo(nil, nil, &mobile)
	if err != nil {
		panic(err)
	}
	if userInfo.Account != "common" {
		panic("wrong account")
	}
	account := "common"
	userInfo, err = model.GetUserInfo(&account, nil, nil)
	if err != nil {
		panic(err)
	}
	if userInfo.Account != "common" {
		panic("wrong account")
	}
}

func TestFindUserInfos(t *testing.T) {
	SetupTestDatabase()
	_ = model.CreateUser("common", "password", model.RoleUser)
	userInfo, err := model.FindUserInfos("%co%", 0, 1)
	if err != nil {
		panic(err)
	}
	if userInfo[0].Account != "common" {
		panic("wrong account")
	}
}

func TestCountUser(t *testing.T) {
	SetupTestDatabase()
	_ = model.CreateUser("common", "password", model.RoleUser)
	count, err := model.CountUser("%co%")
	if err != nil {
		panic(err)
	}
	if count == 0 {
		panic("wrong count")
	}
}
