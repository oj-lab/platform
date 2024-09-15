package problem_service

import (
	"context"

	judge_model "github.com/oj-lab/platform/models/judge"
	problem_model "github.com/oj-lab/platform/models/problem"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
)

func GetProblemInfoList(
	_ context.Context,
	account string,
	titleQuery string,
	limit, offset *int,
) ([]problem_model.Problem, int64, error) {
	db := gorm_agent.GetDefaultDB()
	getOptions := problem_model.GetProblemOptions{
		Selection:  problem_model.ProblemInfoSelection,
		TitleQuery: titleQuery,
		Limit:      limit,
		Offset:     offset,
	}

	total, err := problem_model.CountProblemByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}
	problemList, err := problem_model.GetProblemListByOptions(db, getOptions)
	if err != nil {
		return nil, 0, err
	}

	if len(account) > 0 {
		problemSlugs := []string{}
		for _, problem := range problemList {
			problemSlugs = append(problemSlugs, problem.Slug)
		}
		acceptedJudgeList, err := judge_model.GetJudgeListByOptions(db, judge_model.GetJudgeOptions{
			UserAccount:  account,
			ProblemSlugs: problemSlugs,
			Statuses:     []judge_model.JudgeStatus{judge_model.JudgeStatusFinished},
			Verdicts:     []judge_model.JudgeVerdict{judge_model.JudgeVerdictAccepted},
		})
		if err != nil {
			return nil, 0, err
		}
		acceptedProblemSlugs := map[string]bool{}
		for _, judge := range acceptedJudgeList {
			acceptedProblemSlugs[judge.ProblemSlug] = true
		}
		for i, problem := range problemList {
			if _, ok := acceptedProblemSlugs[problem.Slug]; ok {
				problemList[i].IsAccepted = true
			}
		}
	}

	return problemList, total, nil
}
