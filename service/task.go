package service

import (
	"context"

	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/business"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/google/uuid"
)

func PickJudgeTask(ctx context.Context, consumer string) (*model.JudgeTask, error) {
	task, err := business.GetTaskFromStream(ctx, consumer)
	if err != nil {
		return nil, err
	}

	db := gormAgent.GetDefaultDB()
	err = mapper.UpdateSubmission(db, model.JudgeTaskSubmission{
		UID:    uuid.MustParse(task.SubmissionUID),
		Status: model.SubmissionStatusRunning,
	})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func ReportJudgeTaskResult(
	ctx context.Context,
	consumer string, streamID string, verdictJson string,
) error {
	db := gormAgent.GetDefaultDB()
	err := mapper.UpdateSubmission(db, model.JudgeTaskSubmission{
		RedisStreamID: streamID,
		Status:        model.SubmissionStatusFinished,
		VerdictJson:   verdictJson,
	})
	if err != nil {
		return err
	}

	err = business.AckTaskFromStream(ctx, consumer, streamID)
	if err != nil {
		return err
	}

	return nil
}
