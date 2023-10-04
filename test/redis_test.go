package test

import (
	"context"
	"fmt"
	"testing"

	oj_lab_proto "github.com/OJ-lab/oj-lab-services/service/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
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
	fmt.Println("Get key: ", val)

	pubsub := rdb.Subscribe(ctx, "mychannel")
	ch := pubsub.Channel()

	waitGroup := make(chan struct{})
	go func() {
		fmt.Println("Start listen event...")
		for msg := range ch {
			message := oj_lab_proto.StreamResponse{}
			err := proto.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Received from '%s': %+v\n", msg.Channel, &message)
			close(waitGroup)
		}
	}()

	message := oj_lab_proto.StreamResponse{
		Body: &oj_lab_proto.StreamResponse_Data{
			Data: "hello world",
		},
	}

	data, err := proto.Marshal(&message)
	if err != nil {
		panic(err)
	}

	err = rdb.Publish(ctx, "mychannel", data).Err()
	if err != nil {
		panic(err)
	}

	<-waitGroup
	fmt.Println("Subscriber received message!")

	err = pubsub.Close()
	if err != nil {
		panic(err)
	}
}
