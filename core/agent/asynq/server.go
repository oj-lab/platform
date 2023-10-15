package asynq

import (
	"log"

	"github.com/hibiken/asynq"
)

type AsynqMux struct {
	Pattern string
	*asynq.ServeMux
}

type AsynqServerConfig struct {
	Concurrency int
	UsePriority bool
}

func RunServer(config AsynqServerConfig, subMuxs ...AsynqMux) {
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
	for _, subMux := range subMuxs {
		mux.Handle(subMux.Pattern, subMux.ServeMux)
	}

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
