package auth

import (
	"context"
	"encoding/json"

	redisAgent "github.com/OJ-lab/oj-lab-services/core/agent/redis"
	"github.com/google/uuid"
)

type LoginSession struct {
	Id      string `json:"-"`
	Account string `json:"account"`
}

func NewLoginSession(account string) *LoginSession {
	id := uuid.New().String()

	return &LoginSession{
		Id:      id,
		Account: account,
	}
}

func (ls *LoginSession) SaveToRedis(ctx context.Context) error {
	lsBytes, err := json.Marshal(ls)
	if err != nil {
		return err
	}
	lsString := string(lsBytes)

	err = redisAgent.SetLoginSession(ctx, ls.Id, lsString)
	if err != nil {
		return err
	}

	return nil
}

func CheckLoginSession(ctx context.Context, id string) (*LoginSession, error) {
	lsString, err := redisAgent.GetLoginSession(ctx, id)
	if err != nil {
		return nil, err
	}

	ls := &LoginSession{}
	err = json.Unmarshal([]byte(*lsString), ls)
	if err != nil {
		return nil, err
	}

	return ls, nil
}
