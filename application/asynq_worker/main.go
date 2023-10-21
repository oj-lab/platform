package main

import (
	"github.com/OJ-lab/oj-lab-services/core"
	asynqAgent "github.com/OJ-lab/oj-lab-services/core/agent/asynq"
	"github.com/OJ-lab/oj-lab-services/service/business"
)

func main() {
	core.GetAppLogger().Info("Starting task server...")
	config := asynqAgent.AsynqServerConfig{
		Concurrency: 10,
		UsePriority: true,
	}
	asynqAgent.RunServer(config, business.GetAsynqMuxJudger())
}
