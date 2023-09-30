package business

import (
	"context"
	"encoding/json"
	"fmt"

	asynqAgent "github.com/OJ-lab/oj-lab-services/core/agent/asynq"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const JudgeTaskName = "judge"

var JudgeTaskHandler = asynqAgent.AsynqTaskHandler{
	TaskName:     JudgeTaskName,
	InnerHandler: HandleJudgeTask,
}

func EnqueueJudgeTask(ctx context.Context, judgeTask model.JudgeTask) error {
	client := asynqAgent.GetDefaultTaskClient()
	client.EnqueueTask(JudgeTaskName, judgeTask)

	return nil
}

func HandleJudgeTask(ctx context.Context, task *asynq.Task) error {
	var judgeTask model.JudgeTask
	if err := json.Unmarshal(task.Payload(), &judgeTask); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logrus.Infof("handle judge task: %+v", judgeTask)

	return nil
}
