package judge_model

import (
	"context"

	redis_agent "github.com/oj-lab/platform/modules/agent/redis"
	"github.com/redis/go-redis/v9"
)

const (
	streamName          = "oj_lab_judge_stream"
	consumerGroupName   = "oj_lab_judge_stream_consumer_group"
	defaultConsumerName = "oj_lab_judge_stream_consumer_default"
)

func init() {
	redis_agent := redis_agent.GetDefaultRedisClient()
	_, err := redis_agent.XGroupCreateMkStream(
		context.Background(), streamName, consumerGroupName, "0").Result()
	if err != nil &&
		err != redis.Nil &&
		err.Error() != "BUSYGROUP Consumer Group name already exists" {
		panic(err)
	}
}

func AddTaskToStream(ctx context.Context, task *JudgeTask) (*string, error) {
	redis_agent := redis_agent.GetDefaultRedisClient()
	id, err := redis_agent.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: task.ToStringMap(),
	}).Result()
	if err != nil {
		return nil, err
	}

	return &id, err
}

func GetTaskFromStream(ctx context.Context, consumer string) (*JudgeTask, error) {
	redis_agent := redis_agent.GetDefaultRedisClient()
	if consumer == "" {
		consumer = defaultConsumerName
	}
	result, err := redis_agent.XReadGroup(ctx, &redis.XReadGroupArgs{
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

	task := JudgeTask{}
	for _, message := range result[0].Messages {
		task = *JudgeTaskFromMap(message.Values)
		task.RedisStreamID = &message.ID
	}

	return &task, nil
}

func AckTaskFromStream(ctx context.Context, streamID string) error {
	redis_agent := redis_agent.GetDefaultRedisClient()

	_, err := redis_agent.XAck(ctx, streamName, consumerGroupName, streamID).Result()
	if err != nil {
		return err
	}

	return nil
}
