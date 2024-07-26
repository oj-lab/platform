package auth_module

import (
	"log"
	"testing"
)

func TestGenerateTokenString(t *testing.T) {
	tokenString, err := GenerateAuthTokenString("account", []string{"admin"}...)
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	tokenString, err := GenerateAuthTokenString("account", []string{"admin"}...)
	if err != nil {
		panic(err)
	}
	account, role, err := ParseAuthTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
