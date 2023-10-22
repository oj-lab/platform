package main

import (
	asynqAgent "github.com/OJ-lab/oj-lab-services/core/agent/asynq"
	"github.com/OJ-lab/oj-lab-services/service/business"
)

func main() {
	asynqAgent.RunSecheduler(
		asynqAgent.ScheduleTask{
			Cronspec: "@every 1s",
			Task:     business.NewTaskJudgerTrackAllState(),
		},
	)
}
