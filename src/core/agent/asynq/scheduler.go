package asynqAgent

import (
	"github.com/OJ-lab/oj-lab-services/src/core"
	"github.com/hibiken/asynq"
)

type ScheduleTask struct {
	Cronspec string
	Task     *asynq.Task
}

func RunSecheduler(scheduleTasks ...ScheduleTask) {
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisHost},
		nil,
	)

	for _, scheduleTask := range scheduleTasks {
		entryID, err := scheduler.Register(scheduleTask.Cronspec, scheduleTask.Task)
		if err != nil {
			panic(err)
		}
		core.AppLogger().Infof(
			"setup schedule task: %s, cronspec: %s, entryID: %s",
			scheduleTask.Task.Type(), scheduleTask.Cronspec, entryID,
		)
	}

	scheduler.Run()
}
