package tests

import (
	"github.com/OJ-lab/oj-lab-services/utils"
	"log"
	"testing"
)

func TestGenerateTokenString(t *testing.T) {
	secret := []byte("secret")
	durationString := "1s"
	account := "account"
	tokenString, err := utils.GenerateTokenString(secret, durationString, account)
	if err != nil {
		panic(err)
	}
	log.Print(tokenString)
}

func TestParseTokenString(t *testing.T) {
	secret := []byte("secret")
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiYWNjb3VudCIsImV4cCI6MTY1OTM0MTgxN30.QMN04g45g9NTWZUD3Ys8A0D46BDqMQHEvratR325edU"
	_, err := utils.ParseTokenString(secret, tokenString)
	if err == nil {
		panic("should be expired")
	}
}
