package core_test

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/core/auth"
)

func TestGenerateTokenString(t *testing.T) {
	tokenString, err := auth.GenerateAuthTokenString("account", []string{"admin"}...)
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	tokenString, err := auth.GenerateAuthTokenString("account", []string{"admin"}...)
	if err != nil {
		panic(err)
	}
	account, role, err := auth.ParseAuthTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
