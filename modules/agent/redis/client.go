package redis_agent

import (
	config_module "github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/redis/go-redis/v9"
)

const (
	redisHostsProp = "redis.hosts"
)

var (
	RedisHosts []string
)

func init() {
	RedisHosts = config_module.AppConfig().GetStringSlice(redisHostsProp)
}

type RedisClientInterface interface {
	Set(key string, value interface{}, expiration int64) *redis.StatusCmd
}

var redisClient redis.UniversalClient

func GetDefaultRedisClient() redis.UniversalClient {
	if redisClient == nil {
		if len(RedisHosts) == 0 {
			panic("No redis hosts configured")
		}
		if len(RedisHosts) == 1 {
			redisClient = redis.NewClient(&redis.Options{
				Addr: RedisHosts[0],
			})
		}
		if len(RedisHosts) > 1 {
			redisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: RedisHosts,
			})
		}
	}
	return redisClient
}
