package judge_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

var ErrInvalidJudgeStatus = fmt.Errorf("invalid judge status")

func CreateJudgeResult(
	ctx context.Context,
	judgeResult judge_model.JudgeResult,
) (*judge_model.JudgeResult, error) {
	db := gorm_agent.GetDefaultDB()
	judge, err := GetJudge(ctx, judgeResult.JudgeUID)
	if err != nil {
		return nil, err
	}
	if judge == nil {
		return nil, ErrJudgeNotFound
	}
	if judge.Status != judge_model.JudgeTaskStatusRunning {
		return nil, ErrInvalidJudgeStatus
	}

	newJudgeResult, err := judge_model.CreateJudgeResult(db, judgeResult)
	if err != nil {
		return nil, err
	}
	return newJudgeResult, nil
}

func ReportJudgeResultCount(
	ctx context.Context,
	judgeUID uuid.UUID, resultCount uint,
) error {
	db := gorm_agent.GetDefaultDB()
	return judge_model.UpdateJudge(db, judge_model.Judge{
		UID:         judgeUID,
		ResultCount: resultCount,
	})
}

func ReportJudgeResult(
	ctx context.Context,
	result judge_model.JudgeResult,
) (*judge_model.JudgeResult, error) {
	db := gorm_agent.GetDefaultDB()
	return judge_model.CreateJudgeResult(db, result)
}
