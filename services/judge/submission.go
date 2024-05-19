package judge

import (
	"context"

	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gormAgent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func GetJudgeTaskSubmission(ctx context.Context, uid string) (*judge_model.JudgeTaskSubmission, error) {
	db := gormAgent.GetDefaultDB()
	submission, err := judge_model.GetSubmission(db, uid)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func GetJudgeTaskSubmissionList(
	ctx context.Context, options judge_model.GetSubmissionOptions,
) ([]*judge_model.JudgeTaskSubmission, int64, error) {
	db := gormAgent.GetDefaultDB()
	submissions, total, err := judge_model.GetSubmissionListByOptions(db, options)
	if err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func CreateJudgeTaskSubmission(
	ctx context.Context, submission judge_model.JudgeTaskSubmission,
) (*judge_model.JudgeTaskSubmission, error) {
	db := gormAgent.GetDefaultDB()
	newSubmission, err := judge_model.CreateSubmission(db, submission)
	if err != nil {
		return nil, err
	}

	task := newSubmission.ToJudgeTask()
	streamId, err := addTaskToStream(ctx, &task)
	if err != nil {
		return nil, err
	}

	newSubmission.RedisStreamID = *streamId
	err = judge_model.UpdateSubmission(db, *newSubmission)
	if err != nil {
		return nil, err
	}

	return newSubmission, nil
}
