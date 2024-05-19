package judge

import (
	"context"

	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	redis_agent "github.com/oj-lab/oj-lab-platform/modules/agent/redis"
	"github.com/redis/go-redis/v9"
)

const (
	streamName          = "oj_lab_judge_stream"
	consumerGroupName   = "oj_lab_judge_stream_consumer_group"
	defaultConsumerName = "oj_lab_judge_stream_consumer_default"
)

func init() {
	redisAgent := redis_agent.GetDefaultRedisClient()
	_, err := redisAgent.XGroupCreateMkStream(context.Background(), streamName, consumerGroupName, "0").Result()
	if err != nil && err != redis.Nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		panic(err)
	}
}

func addTaskToStream(ctx context.Context, task *judge_model.JudgeTask) (*string, error) {
	redisAgent := redis_agent.GetDefaultRedisClient()
	id, err := redisAgent.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: task.ToStringMap(),
	}).Result()
	if err != nil {
		return nil, err
	}

	return &id, err
}

func getTaskFromStream(ctx context.Context, consumer string) (*judge_model.JudgeTask, error) {
	redisAgent := redis_agent.GetDefaultRedisClient()
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

	task := judge_model.JudgeTask{}
	for _, message := range result[0].Messages {
		task = *judge_model.JudgeTaskFromMap(message.Values)
		task.RedisStreamID = &message.ID
	}

	return &task, nil
}

func ackTaskFromStream(ctx context.Context, consumer string, streamID string) error {
	redisAgent := redis_agent.GetDefaultRedisClient()
	// TODO: Some ineffectual assignment here, need to find out why
	// if consumer == "" {
	// 	consumer = defaultConsumerName
	// }
	_, err := redisAgent.XAck(ctx, streamName, consumerGroupName, streamID).Result()
	if err != nil {
		return err
	}

	return nil
}
