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

	total, err := problem_model.CountProblemByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}
	problemList, err := problem_model.GetProblemListByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	return problemList, total, nil
}
