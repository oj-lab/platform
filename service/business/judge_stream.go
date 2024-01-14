package business

import (
	"context"

	redisAgent "github.com/OJ-lab/oj-lab-services/core/agent/redis"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/redis/go-redis/v9"
)

const (
	streamName          = "oj_lab_judge_stream"
	consumerGroupName   = "oj_lab_judge_stream_consumer_group"
	defaultConsumerName = "oj_lab_judge_stream_consumer_default"
)

func init() {
	redisAgent := redisAgent.GetDefaultRedisClient()
	_, err := redisAgent.XGroupCreateMkStream(context.Background(), streamName, consumerGroupName, "0").Result()
	if err != nil && err != redis.Nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		panic(err)
	}
}

func AddTaskToStream(ctx context.Context, task *model.JudgeTask) (*string, error) {
	redisAgent := redisAgent.GetDefaultRedisClient()
	id, err := redisAgent.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: task.ToStringMap(),
	}).Result()
	if err != nil {
		return nil, err
	}

	return &id, err
}

func GetTaskFromStream(ctx context.Context, consumer string) (*model.JudgeTask, error) {
	redisAgent := redisAgent.GetDefaultRedisClient()
	if consumer == "" {
		consumer = defaultConsumerName
	}
	result, err := redisAgent.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroupName,
		Consumer: consumer,
		Streams:  []string{streamName, ">"},
		Count:    1,
		Block:    -1,
	}).Result()

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	task := model.JudgeTask{}
	for _, message := range result[0].Messages {
		task = *model.JudgeTaskFromMap(message.Values)
		task.RedisStreamID = &message.ID
	}

	return &task, nil
}

func AckTaskFromStream(ctx context.Context, consumer string, streamID string) error {
	redisAgent := redisAgent.GetDefaultRedisClient()
	if consumer == "" {
		consumer = defaultConsumerName
	}
	_, err := redisAgent.XAck(ctx, streamName, consumerGroupName, streamID).Result()
	if err != nil {
		return err
	}

	return nil
}
