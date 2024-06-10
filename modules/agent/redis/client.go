package redisAgent

import (
	"github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/redis/go-redis/v9"
)

const (
	redisHostsProp = "redis.hosts"
)

var (
	redisHosts []string
)

func init() {
	redisHosts = config.AppConfig.GetStringSlice(redisHostsProp)
}

type RedisClientInterface interface {
	Set(key string, value interface{}, expiration int64) *redis.StatusCmd
}

var redisClient redis.UniversalClient

func GetDefaultRedisClient() redis.UniversalClient {
	if redisClient == nil {
		if len(redisHosts) == 0 {
			panic("No redis hosts configured")
		}
		if len(redisHosts) == 1 {
			redisClient = redis.NewClient(&redis.Options{
				Addr: redisHosts[0],
			})
		}
		if len(redisHosts) > 1 {
			redisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: redisHosts,
			})
		}
	}
	return redisClient
}
