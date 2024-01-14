package business

import (
	"context"

	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
)

func GetProblemInfoList(ctx context.Context) ([]model.ProblemInfo, int64, error) {
	db := gormAgent.GetDefaultDB()
	getOptions := mapper.GetProblemOptions{
		Selection: model.ProblemInfoSelection,
	}

	problemList, total, err := mapper.GetProblemListByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	problemInfoList := []model.ProblemInfo{}
	for _, problem := range problemList {
		problemInfoList = append(problemInfoList, problem.ToProblemInfo())
	}
	return problemInfoList, total, nil
}
