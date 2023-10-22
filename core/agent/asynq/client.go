package asynqAgent

import (
	"github.com/hibiken/asynq"
)

type AsynqClient struct {
	innerClient *asynq.Client
}

var client *AsynqClient

func GetDefaultTaskClient() *AsynqClient {
	if client == nil {
		client = &AsynqClient{
			innerClient: asynq.NewClient(asynq.RedisClientOpt{
				Addr: redisHost,
			}),
		}
	}

	return client
}

func (ac *AsynqClient) EnqueueTask(task *asynq.Task, opts ...asynq.Option) error {
	_, err := ac.innerClient.Enqueue(task)
	return err
}
