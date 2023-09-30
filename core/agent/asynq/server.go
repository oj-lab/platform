package asynq

import (
	"log"

	"github.com/hibiken/asynq"
)

type AsynqTaskHandler struct {
	TaskName     string
	InnerHandler asynq.HandlerFunc
}

type AsynqServerConfig struct {
	Concurrency int
	UsePriority bool
}

func RunServer(config AsynqServerConfig, handlers ...AsynqTaskHandler) {
	asynqConfig := asynq.Config{
		Concurrency: config.Concurrency,
	}
	if config.UsePriority {
		asynqConfig.Queues = priorityMap
	}

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisHost},
		asynqConfig,
	)

	mux := asynq.NewServeMux()
	for _, handler := range handlers {
		mux.HandleFunc(handler.TaskName, handler.InnerHandler)
	}

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
