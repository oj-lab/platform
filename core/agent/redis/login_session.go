package redisAgent

import (
	"context"
	"fmt"
	"time"
)

const loginSessionKeyFormat = "LS_%s"
const loginSessionDuration = time.Second * 30

func SetLoginSession(ctx context.Context, sessionId string, sessionString string) error {
	redisClient := GetDefaultRedisClient()
	key := fmt.Sprintf(loginSessionKeyFormat, sessionId)

	err := redisClient.Set(ctx, key, sessionString, loginSessionDuration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetLoginSession(ctx context.Context, sessionId string) (*string, error) {
	redisClient := GetDefaultRedisClient()
	key := fmt.Sprintf(loginSessionKeyFormat, sessionId)

	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &val, nil
}
