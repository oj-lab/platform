package core_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	oj_lab_proto "github.com/OJ-lab/oj-lab-services/src/service/proto"
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
		defer close(waitGroup)
		fmt.Println("Start listen event...")
		timeout := time.After(5 * time.Second)
		for {
			select {
			case msg := <-ch:
				message := oj_lab_proto.StreamResponse{}
				err := proto.Unmarshal([]byte(msg.Payload), &message)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Received from '%s': %+v\n", msg.Channel, &message)
				return
			case <-timeout:
				panic("timeout")
			}
		}
	}()

	fmt.Println("Wait for 1 second...")
	time.Sleep(1 * time.Second)

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
	fmt.Println("Published message...")

	<-waitGroup
	fmt.Println("Subscriber received message!")

	err = pubsub.Close()
	if err != nil {
		panic(err)
	}

	sId := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "oj-lab-judge-test",
		Values: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}).Val()
	println("sId: ", sId)
	rdb.XGroupCreate(ctx, "oj-lab-judge-test", "oj-lab-judge-test-group", "0")
	ack := rdb.XAck(ctx, "oj-lab-judge-test", "oj-lab-judge-test-group", sId).Val()
	if ack == 1 {
		panic("should not ack")
	}
	rdb.XGroupCreateConsumer(ctx, "oj-lab-judge-test", "oj-lab-judge-test-group", "oj-lab-judge-test-consumer-1")
	sSlice := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "oj-lab-judge-test-group",
		Consumer: "oj-lab-judge-test-consumer-1",
		Streams:  []string{"oj-lab-judge-test", ">"},
		Count:    1,
		Block:    0,
	}).Val()
	sVal := sSlice[0].Messages[0].Values
	for k, v := range sVal {
		fmt.Printf("k: %s, v: %s\n", k, v)
	}
	ack = rdb.XAck(ctx, "oj-lab-judge-test", "oj-lab-judge-test-group", sSlice[0].Messages[0].ID).Val()
	if ack != 1 {
		panic("failed to ack")
	}
}
