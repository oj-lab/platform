package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func TestUserDB(t *testing.T) {
	db := gorm_agent.GetDefaultDB()
	user := user_model.User{
		Account:  "test",
		Password: func() *string { s := "test"; return &s }(),
	}
	err := user_model.CreateUser(db, user)
	if err != nil {
		t.Error(err)
	}

	dbUser, err := user_model.GetUser(db, user.Account)
	if err != nil {
		t.Error(err)
	}
	userJson, err := json.MarshalIndent(dbUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(userJson))

	dbPublicUser, err := user_model.GetPublicUser(db, user.Account)
	if err != nil {
		t.Error(err)
	}
	publicUserJson, err := json.MarshalIndent(dbPublicUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(publicUserJson))

	err = user_model.DeleteUser(db, user)
	if err != nil {
		t.Error(err)
	}
}
