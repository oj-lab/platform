package judge_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

var ErrJudgeNotFound = fmt.Errorf("judge not found")

func GetJudge(ctx context.Context, uid uuid.UUID) (*judge_model.Judge, error) {
	db := gorm_agent.GetDefaultDB()
	judge, err := judge_model.GetJudge(db, uid)
	if err != nil {
		return nil, err
	}
	return judge, nil
}

func GetJudgeList(
	ctx context.Context, options judge_model.GetJudgeOptions,
) ([]*judge_model.Judge, int64, error) {
	db := gorm_agent.GetDefaultDB()
	judges, total, err := judge_model.GetJudgeListByOptions(db, options)
	if err != nil {
		return nil, 0, err
	}

	return judges, total, nil
}

func CreateJudge(
	ctx context.Context, judge judge_model.Judge,
) (*judge_model.Judge, error) {
	db := gorm_agent.GetDefaultDB()
	newJudge, err := judge_model.CreateJudge(db, judge)
	if err != nil {
		return nil, err
	}

	task := newJudge.ToJudgeTask()
	streamId, err := judge_model.AddTaskToStream(ctx, &task)
	if err != nil {
		return nil, err
	}

	newJudge.RedisStreamID = *streamId
	err = judge_model.UpdateJudge(db, *newJudge)
	if err != nil {
		return nil, err
	}

	return newJudge, nil
}
