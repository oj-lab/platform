package auth

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type LoginSession struct {
	Id      uuid.UUID `json:"-"`
	Account string    `json:"account"`
	Roles   []string  `json:"roles"`
}

func (ls LoginSession) GetJsonString() (string, error) {
	lsBytes, err := json.Marshal(ls)
	if err != nil {
		return "", err
	}
	lsString := string(lsBytes)

	return lsString, nil
}

func GetLoginSessionFromJsonString(lsString string) (*LoginSession, error) {
	ls := &LoginSession{}
	err := json.Unmarshal([]byte(lsString), ls)
	if err != nil {
		return nil, err
	}
	return ls, nil
}

func NewLoginSession(ls LoginSession) *LoginSession {
	return &LoginSession{
		Id:      uuid.New(),
		Account: ls.Account,
		Roles:   ls.Roles,
	}
}

func (ls LoginSession) SaveToRedis(ctx context.Context) error {
	err := SetLoginSession(ctx, ls)
	if err != nil {
		return err
	}

	return nil
}
