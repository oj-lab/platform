package service

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/core"
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func GetJudgeTaskSubmissionList(
	ctx context.Context, options mapper.GetSubmissionOptions,
) ([]*model.JudgeTaskSubmission, int64, *core.SeviceError) {
	db := gormAgent.GetDefaultDB()
	submissions, total, err := mapper.GetSubmissionListByOptions(db, options)
	if err != nil {
		return nil, 0, core.NewInternalError("failed to get submission list")
	}

	return submissions, total, nil
}

func CreateJudgeTaskSubmission(
	ctx context.Context, submission model.JudgeTaskSubmission,
) (*model.JudgeTaskSubmission, *core.SeviceError) {
	db := gormAgent.GetDefaultDB()
	newSubmission, err := mapper.CreateSubmission(db, submission)
	if err != nil {
		return nil, core.NewInternalError("failed to create submission")
	}

	return newSubmission, nil
}
