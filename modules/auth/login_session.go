package auth_module

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	config_module "github.com/oj-lab/platform/modules/config"
)

const defaultLoginSessionDuration = time.Hour * 24 * 7

var LoginSessionDuration time.Duration

func init() {
	LoginSessionDuration = config_module.AppConfig().GetDuration("service.login_session_duration")
	if LoginSessionDuration == 0 {
		LoginSessionDuration = defaultLoginSessionDuration
	}
}

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
