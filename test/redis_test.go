package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func TestRedis(t *testing.T) {
	// rdb := redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},

	// 	// To route commands by latency or randomly, enable one of the following.
	// 	// RouteByLatency: true,
	// 	// RouteRandomly: true,
	// })

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
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

	// 订阅频道
	pubsub := rdb.Subscribe(ctx, "mychannel")

	// 从通道中读取消息
	ch := pubsub.Channel()

	waitGroup := make(chan struct{})
	// 在 goroutine 中处理消息
	go func() {
		fmt.Println("start")
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload)
			close(waitGroup)
		}
	}()

	// 发布消息
	err = rdb.Publish(ctx, "mychannel", "hello world").Err()
	if err != nil {
		panic(err)
	}

	<-waitGroup

	// 关闭订阅
	err = pubsub.Close()
	if err != nil {
		panic(err)
	}
}
