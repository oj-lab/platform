package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSettings *JWTSettings

func SetupJWTSettings(settings JWTSettings) {
	jwtSettings = &settings
}

func GenerateTokenString(account string, role string) (string, error) {
	duration, err := time.ParseDuration(jwtSettings.Duration)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account,
		"role":    role,
		"exp":     time.Now().Add(duration).Unix(),
	})
	return token.SignedString([]byte(jwtSettings.Secret))
}

func ParseTokenString(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSettings.Secret), nil
	})
	if err != nil {
		return "", "", err
	}
	if token.Valid {
		return token.Claims.(jwt.MapClaims)["account"].(string), token.Claims.(jwt.MapClaims)["role"].(string), nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return "", "", errors.New("token is expired")
		}
	}
	return "", "", errors.New("invalid token")
}
