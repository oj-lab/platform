package main

import (
	"github.com/OJ-lab/oj-lab-services/src/core"
	asynqAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/asynq"
	"github.com/OJ-lab/oj-lab-services/src/service/business"
)

func main() {
	core.AppLogger().Info("Starting task server...")
	config := asynqAgent.AsynqServerConfig{
		Concurrency: 10,
		UsePriority: true,
	}
	asynqAgent.RunServer(config, business.GetAsynqMuxJudger())
}
