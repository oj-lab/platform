package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateTokenString(secret []byte, durationString string, account string) (string, error) {
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account,
		"exp":     time.Now().Add(duration).Unix(),
	})
	return token.SignedString(secret)
}

func ParseTokenString(secret []byte, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if token.Valid {
		return token.Claims.(jwt.MapClaims)["account"].(string), nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return "", errors.New("token is expired")
		}
	}
	return "", errors.New("invalid token")
}
