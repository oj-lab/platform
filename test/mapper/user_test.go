package mapper_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/package/mapper"
	"github.com/OJ-lab/oj-lab-services/package/model"
)

func TestUserMapper(t *testing.T) {
	password := "test"
	user := model.User{
		Account:  "test",
		Password: &password,
		Roles:    []*model.Role{{Name: "admin"}},
	}
	err := mapper.CreateUser(user)
	if err != nil {
		t.Error(err)
	}

	dbUser, err := mapper.GetUser(user.Account)
	if err != nil {
		t.Error(err)
	}

	userJson, err := json.MarshalIndent(dbUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(userJson))
}
