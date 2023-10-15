package main

import (
	asynqAgent "github.com/OJ-lab/oj-lab-services/core/agent/asynq"
	"github.com/OJ-lab/oj-lab-services/service/business"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting task server...")
	config := asynqAgent.AsynqServerConfig{
		Concurrency: 10,
		UsePriority: true,
	}
	asynqAgent.RunServer(config, business.GetJudgeMux())
}
