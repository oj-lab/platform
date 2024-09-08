package auth_module

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	config_module "github.com/oj-lab/platform/modules/config"
)

var (
	jwtSecret   string
	jwtDuration time.Duration
)

func init() {
	jwtSecret = config_module.AppConfig().GetString("jwt.secret")
	jwtDuration = config_module.AppConfig().GetDuration("jwt.duration")
}

type AuthToken struct {
	Account string
	Roles   []string
	Expires time.Time
}

func (a *AuthToken) Valid() error {
	if a.Expires.Before(time.Now()) {
		return errors.New("token is expired")
	}
	return nil
}

func (a *AuthToken) ToJWTMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"account": a.Account,
		"roles":   a.Roles,
		"exp":     a.Expires.Unix(),
	}
}

func GetAuthTokenFromJWTMapClaims(claims jwt.MapClaims) *AuthToken {
	roleInterface := claims["roles"].([]interface{})
	roles := make([]string, len(roleInterface))
	for i, role := range roleInterface {
		roles[i] = role.(string)
	}
	return &AuthToken{
		Account: claims["account"].(string),
		Roles:   roles,
		Expires: time.Unix(int64(claims["exp"].(float64)), 0),
	}
}

func GenerateAuthTokenString(account string, roles ...string) (string, error) {
	duration := jwtDuration
	token := AuthToken{
		Account: account,
		Roles:   roles,
		Expires: time.Now().Add(duration),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.ToJWTMapClaims())
	return jwtToken.SignedString([]byte(jwtSecret))
}

func ParseAuthTokenString(tokenString string) (string, []string, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if jwtToken.Valid {
		claims := jwtToken.Claims.(jwt.MapClaims)
		token := GetAuthTokenFromJWTMapClaims(claims)
		return token.Account, token.Roles, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return "", nil, errors.New("token is expired")
		}
	}
	return "", nil, errors.New("invalid token")
}
