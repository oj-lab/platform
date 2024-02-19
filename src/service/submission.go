package service

import (
	"context"
	"fmt"

	"github.com/OJ-lab/oj-lab-services/src/core"
	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/src/service/business"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
)

func GetJudgeTaskSubmission(ctx context.Context, uid string) (*model.JudgeTaskSubmission, *core.SeviceError) {
	db := gormAgent.GetDefaultDB()
	submission, err := mapper.GetSubmission(db, uid)
	if err != nil {
		return nil, core.NewInternalError("failed to get submission by uid")
	}

	return submission, nil
}


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

	task := newSubmission.ToJudgeTask()
	streamId, err := business.AddTaskToStream(ctx, &task)
	if err != nil {
		return nil, core.NewInternalError(fmt.Sprintf("failed to add task to stream %v", err))
	}

	newSubmission.RedisStreamID = *streamId
	err = mapper.UpdateSubmission(db, *newSubmission)
	if err != nil {
		return nil, core.NewInternalError("failed to update submission")
	}

	return newSubmission, nil
}
