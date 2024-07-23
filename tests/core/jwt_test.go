package core_test

import (
	"log"
	"testing"

	auth_module "github.com/oj-lab/oj-lab-platform/modules/auth"
)

func TestGenerateTokenString(t *testing.T) {
	tokenString, err := auth_module.GenerateAuthTokenString("account", []string{"admin"}...)
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	tokenString, err := auth_module.GenerateAuthTokenString("account", []string{"admin"}...)
	if err != nil {
		panic(err)
	}
	account, role, err := auth_module.ParseAuthTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
