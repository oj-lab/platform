package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	redisAgent "github.com/oj-lab/oj-lab-platform/modules/agent/redis"
)

const loginSessionKeyFormat = "LS_%s"
const loginSessionDuration = time.Second * 30

func SetLoginSession(ctx context.Context, ls LoginSession) error {
	redisClient := redisAgent.GetDefaultRedisClient()
	key := fmt.Sprintf(loginSessionKeyFormat, ls.Id.String())
	value, err := ls.GetJsonString()
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, key, value, loginSessionDuration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetLoginSession(ctx context.Context, id uuid.UUID) (*LoginSession, error) {
	redisClient := redisAgent.GetDefaultRedisClient()
	lsIdString := id.String()
	key := fmt.Sprintf(loginSessionKeyFormat, lsIdString)

	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	ls, err := GetLoginSessionFromJsonString(val)
	if err != nil {
		return nil, err
	}
	ls.Id = id

	return ls, nil
}

func UpdateLoginSession(ctx context.Context, idString, sesionString string) error {
	redisClient := redisAgent.GetDefaultRedisClient()
	key := fmt.Sprintf(loginSessionKeyFormat, idString)

	err := redisClient.Set(ctx, key, sesionString, loginSessionDuration).Err()
	if err != nil {
		return err
	}

	return nil
}
