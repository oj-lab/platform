package problem_service

import (
	"context"

	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func GetProblemInfoList(_ context.Context, limit, offset *int) ([]problem_model.Problem, int64, error) {
	db := gorm_agent.GetDefaultDB()
	getOptions := problem_model.GetProblemOptions{
		Selection: problem_model.ProblemInfoSelection,
		Limit:     limit,
		Offset:    offset,
	}

	problemList, total, err := problem_model.GetProblemInfoListByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	return problemList, total, nil
}
