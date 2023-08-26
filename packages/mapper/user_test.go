package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/packages/model"
)

func TestUserMapper(t *testing.T) {
	ctx := context.Background()

	password := "test"
	user := model.User{
		Account:  "test",
		Password: &password,
		Roles:    []*model.Role{{Name: "admin"}},
	}
	err := CreateUser(ctx, user)
	if err != nil {
		t.Error(err)
	}

	dbUser, err := GetUser(ctx, user.Account)
	if err != nil {
		t.Error(err)
	}

	userJson, err := json.MarshalIndent(dbUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(userJson))
}
