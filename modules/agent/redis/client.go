package redisAgent

import (
	"github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/redis/go-redis/v9"
)

const (
	redisHostProp = "redis.host"
)

var (
	redisHost string
)

func init() {
	redisHost = config.AppConfig.GetString(redisHostProp)
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
