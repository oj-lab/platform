package auth

import (
	"context"
	"fmt"
	"time"

	redisAgent "github.com/oj-lab/oj-lab-platform/modules/agent/redis"
	"github.com/oj-lab/oj-lab-platform/modules/log"
	"github.com/redis/go-redis/v9"
)

const loginSessionKeyFormat = "LS_%s_%s" // "LS_<account>_<uuid>"
const loginSessionDuration = time.Minute * 15

func getLoginSessionRedisKey(key LoginSessionKey) string {
	return fmt.Sprintf(loginSessionKeyFormat, key.Account, key.Id.String())
}

func SetLoginSession(ctx context.Context, key LoginSessionKey, data LoginSessionData) error {
	redisClient := redisAgent.GetDefaultRedisClient()

	value, err := data.GetJsonString()
	if err != nil {
		return err
	}
	// TODO: Watch Redis JSON SET usage, currently not support atomic SETEX
	err = redisClient.Set(ctx, getLoginSessionRedisKey(key), value, loginSessionDuration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetLoginSession(ctx context.Context, key LoginSessionKey) (*LoginSession, error) {
	redisClient := redisAgent.GetDefaultRedisClient()

	val, err := redisClient.Get(ctx, getLoginSessionRedisKey(key)).Result()
	if err != nil {
		return nil, err
	}
	data, err := getLoginSessionDataFromJsonString(val)
	if err != nil {
		return nil, err
	}

	return &LoginSession{
		Key:  key,
		Data: *data,
	}, nil
}

func UpdateLoginSessionByAccount(ctx context.Context, account string, data LoginSessionData) error {
	redisClient := redisAgent.GetDefaultRedisClient()

	redisKeys, err := redisClient.Keys(ctx, fmt.Sprintf(loginSessionKeyFormat, account, "*")).Result()
	if err != nil {
		return err
	}

	val, err := data.GetJsonString()
	if err != nil {
		return err
	}
	for _, redisKey := range redisKeys {
		// TODO: KeepTTL only works in redis v6+
		err = redisClient.Set(ctx, redisKey, val, redis.KeepTTL).Err()
		if err != nil {
			log.AppLogger().Errorf("failed to update login session: %v", err)
		}
	}

	return nil
}
