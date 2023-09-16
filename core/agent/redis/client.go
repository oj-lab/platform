package redis

import "github.com/redis/go-redis/v9"

var redisClient *redis.Client

func GetDefaultRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}
	return redisClient
}
