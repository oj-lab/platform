package mapper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/packages/model"
)

func TestUserMapper(t *testing.T) {
	password := "test"
	user := model.User{
		Account:  "test",
		Password: &password,
		Roles:    []*model.Role{{Name: "admin"}},
	}
	err := CreateUser(user)
	if err != nil {
		t.Error(err)
	}

	dbUser, err := GetUser(user.Account)
	if err != nil {
		t.Error(err)
	}

	userJson, err := json.MarshalIndent(dbUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(userJson))
}
