package redisAgent

import (
	"github.com/OJ-lab/oj-lab-services/src/core"
	"github.com/redis/go-redis/v9"
)

const (
	redisHostProp = "redis.host"
)

var (
	redisHost string
)

func init() {
	redisHost = core.AppConfig.GetString(redisHostProp)
}

var redisClient *redis.Client

func GetDefaultRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr: redisHost,
		})
	}
	return redisClient
}
