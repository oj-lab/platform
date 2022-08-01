package tests

import (
	"github.com/OJ-lab/oj-lab-services/utils"
	"log"
	"testing"
)

func TestGenerateTokenString(t *testing.T) {
	jwtSettings, _ := utils.GetJWTSettings("../../config/test.ini")
	utils.SetupJWTSettings(jwtSettings)
	tokenString, err := utils.GenerateTokenString("account", "admin")
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	jwtSettings, _ := utils.GetJWTSettings("../../config/test.ini")
	utils.SetupJWTSettings(jwtSettings)
	tokenString, err := utils.GenerateTokenString("account", "admin")
	if err != nil {
		panic(err)
	}
	account, role, err := utils.ParseTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
