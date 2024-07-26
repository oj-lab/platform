package user_model

import (
	"encoding/json"
	"fmt"
	"testing"

	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func TestUserDB(t *testing.T) {
	db := gorm_agent.GetDefaultDB()
	user := User{
		Account:  "test",
		Password: func() *string { s := "test"; return &s }(),
	}
	err := CreateUser(db, user)
	if err != nil {
		t.Error(err)
	}

	dbUser, err := GetUser(db, user.Account)
	if err != nil {
		t.Error(err)
	}
	userJson, err := json.MarshalIndent(dbUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(userJson))

	dbPublicUser, err := GetPublicUser(db, user.Account)
	if err != nil {
		t.Error(err)
	}
	publicUserJson, err := json.MarshalIndent(dbPublicUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(publicUserJson))

	err = DeleteUser(db, user)
	if err != nil {
		t.Error(err)
	}
}
