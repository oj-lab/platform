package business

import (
	"errors"
	"time"

	"github.com/OJ-lab/oj-lab-services/packages/application"
	"github.com/OJ-lab/oj-lab-services/packages/model"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret string
var jwtDuration time.Duration

func init() {
	jwtSecret = application.AppConfig.GetString("jwt.secret")
	jwtDuration = application.AppConfig.GetDuration("jwt.duration")
}

func GenerateTokenString(account string, roles []model.Role) (string, error) {
	duration := jwtDuration

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account,
		"roles":   roles,
		"exp":     time.Now().Add(duration).Unix(),
	})
	return token.SignedString([]byte(jwtSecret))
}

func ParseTokenString(tokenString string) (string, []string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", nil, err
	}
	if token.Valid {
		roleInterface := token.Claims.(jwt.MapClaims)["roles"].([]interface{})
		roles := make([]string, len(roleInterface))
		for i, v := range roleInterface {
			roles[i] = v.(string)
		}
		return token.Claims.(jwt.MapClaims)["account"].(string), roles, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return "", nil, errors.New("token is expired")
		}
	}
	return "", nil, errors.New("invalid token")
}
