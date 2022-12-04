package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func TestRedisCluster(t *testing.T) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"127.0.0.1:7000", "127.0.0.:7001", "127.0.0.:7002", "127.0.0.:7003", "127.0.0.:7004", "127.0.0.:7005"},

		// To route commands by latency or randomly, enable one of the following.
		// RouteByLatency: true,
		// RouteRandomly: true,
	})
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
