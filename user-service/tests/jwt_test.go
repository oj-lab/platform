package tests

import (
	"log"
	"testing"

	"github.com/OJ-lab/oj-lab-services/config"
	"github.com/OJ-lab/oj-lab-services/model"
	"github.com/OJ-lab/oj-lab-services/user-service/business"
)

func TestGenerateTokenString(t *testing.T) {
	jwtSettings, _ := config.GetJWTSettings("../../config/ini/test.ini")
	business.SetupJWTSettings(jwtSettings)
	tokenString, err := business.GenerateTokenString("account", []model.Role{model.RoleAdmin})
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	jwtSettings, _ := config.GetJWTSettings("../../config/ini/test.ini")
	business.SetupJWTSettings(jwtSettings)
	tokenString, err := business.GenerateTokenString("account", []model.Role{model.RoleAdmin})
	if err != nil {
		panic(err)
	}
	account, role, err := business.ParseTokenString(tokenString)
	if err != nil {
		panic(err)
	}
	log.Println(account, role)
}
