package tests

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
)

func TestGenerateTokenString(t *testing.T) {
	tokenString, err := business.GenerateTokenString("account", []*model.Role{{Name: "admin"}})
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	tokenString, err := business.GenerateTokenString("account", []*model.Role{{Name: "admin"}})
	if err != nil {
		panic(err)
	}
	account, role, err := business.ParseTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
