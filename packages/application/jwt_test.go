package application

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/packages/model"
)

func TestGenerateTokenString(t *testing.T) {
	tokenString, err := GenerateTokenString("account", []*model.Role{{Name: "admin"}})
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	tokenString, err := GenerateTokenString("account", []*model.Role{{Name: "admin"}})
	if err != nil {
		panic(err)
	}
	account, role, err := ParseTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
