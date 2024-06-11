package problem_service

import (
	"context"

	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func getProblemInfoList(_ context.Context) ([]problem_model.ProblemInfo, int64, error) {
	db := gorm_agent.GetDefaultDB()
	getOptions := problem_model.GetProblemOptions{
		Selection: problem_model.ProblemInfoSelection,
	}

	problemList, total, err := problem_model.GetProblemListByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	problemInfoList := []problem_model.ProblemInfo{}
	for _, problem := range problemList {
		problemInfoList = append(problemInfoList, problem.ToProblemInfo())
	}
	return problemInfoList, total, nil
}
