package business

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func GetProblemInfoList(ctx context.Context) ([]model.ProblemInfo, int64, error) {
	getOptions := mapper.GetProblemOptions{
		Selection: model.ProblemInfoSelection,
	}

	problemList, total, err := mapper.GetProblemListByOptions(getOptions)
	if err != nil {
		return nil, 0, err
	}

	problemInfoList := []model.ProblemInfo{}
	for _, problem := range problemList {
		problemInfoList = append(problemInfoList, problem.ToProblemInfo())
	}
	return problemInfoList, total, nil
}
