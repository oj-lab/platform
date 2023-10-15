package business

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	asynqAgent "github.com/OJ-lab/oj-lab-services/core/agent/asynq"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type judgeTaskCtxKey string

const (
	WaitJudgerTaskName                           = "judge:wait_judger"
	JudgerPickSubmissionTaskName                 = "judge:judger_pick_submission"
	JudgeTaskName                                = "judge"
	JudgeMuxPattern                              = "judge"
	ctxJudgerHostKey             judgeTaskCtxKey = "judgerHost"
)

var JudgeTaskHandler = asynqAgent.AsynqMux{
	Pattern: JudgeMuxPattern,
}

func GetJudgeMux() asynqAgent.AsynqMux {
	serveMux := asynq.NewServeMux()
	serveMux.HandleFunc(JudgeTaskName, handleJudgeTask)

	return asynqAgent.AsynqMux{
		Pattern:  JudgeMuxPattern,
		ServeMux: serveMux,
	}
}

func EnqueueWaitJudgerTask(ctx context.Context, judger model.Judger) error {
	client := asynqAgent.GetDefaultTaskClient()
	return client.EnqueueTask(WaitJudgerTaskName, judger, asynq.TaskID(judger.Host))
}

func handleWaitJudgerTask(ctx context.Context, task *asynq.Task) error {
	judgerHost := ctx.Value(ctxJudgerHostKey)
	logrus.Infof("judgerHost: %v", judgerHost)

	var judger model.Judger
	if err := json.Unmarshal(task.Payload(), &judger); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	time.Sleep(5 * time.Second)
	logrus.Infof("handle wait judger task: %+v", judger)

	return nil
}

func EnqueueJudgerPickSubmissionTask(ctx context.Context, judger model.Judger) error {
	client := asynqAgent.GetDefaultTaskClient()
	return client.EnqueueTask(JudgerPickSubmissionTaskName, judger, asynq.TaskID(judger.Host))
}

func EnqueueJudgeTask(ctx context.Context, judgeTask model.JudgeTask) error {
	client := asynqAgent.GetDefaultTaskClient()
	return client.EnqueueTask(JudgeTaskName, judgeTask)
}

func handleJudgeTask(ctx context.Context, task *asynq.Task) error {
	judgerHost := ctx.Value(ctxJudgerHostKey)
	logrus.Infof("judgerHost: %v", judgerHost)

	var judgeTask model.JudgeTask
	if err := json.Unmarshal(task.Payload(), &judgeTask); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	time.Sleep(5 * time.Second)
	logrus.Infof("handle judge task: %+v", judgeTask)

	return nil
}
