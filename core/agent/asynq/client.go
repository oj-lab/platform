package asynq

import (
	"encoding/json"

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

func (ac *AsynqClient) EnqueueTask(name string, payload interface{}, opts ...asynq.Option) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(name, payloadBytes, opts...)
	_, err = ac.innerClient.Enqueue(task)

	return err
}
