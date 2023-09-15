package core_test

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/package/core"
	"github.com/OJ-lab/oj-lab-services/package/model"
)

func TestGenerateTokenString(t *testing.T) {
	tokenString, err := core.GenerateTokenString("account", []*model.Role{{Name: "admin"}})
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	tokenString, err := core.GenerateTokenString("account", []*model.Role{{Name: "admin"}})
	if err != nil {
		panic(err)
	}
	account, role, err := core.ParseTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
