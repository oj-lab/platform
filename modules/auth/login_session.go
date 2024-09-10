package auth_module

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type LoginSessionKey struct {
	Account string
	Id      uuid.UUID
}

type LoginSessionData struct {
}

type LoginSession struct {
	Key  LoginSessionKey
	Data LoginSessionData
}

func (data LoginSessionData) GetJsonString() (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataString := string(bytes)

	return dataString, nil
}

func getLoginSessionDataFromJsonString(dataString string) (*LoginSessionData, error) {
	data := &LoginSessionData{}
	err := json.Unmarshal([]byte(dataString), data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewLoginSession(account string, data LoginSessionData) *LoginSession {
	return &LoginSession{
		LoginSessionKey{
			Account: account,
			Id:      uuid.New(),
		},
		LoginSessionData{},
	}
}

func (ls LoginSession) SaveToRedis(ctx context.Context) error {
	err := SetLoginSession(ctx, ls.Key, ls.Data)
	if err != nil {
		return err
	}

	return nil
}

func (ls LoginSession) DelInRedis(ctx context.Context) error {
	err := DelLoginSession(ctx, ls.Key)
	if err != nil {
		return err
	}

	return nil
}
