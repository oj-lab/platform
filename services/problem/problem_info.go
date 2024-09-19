package problem_service

import (
	"context"

	judge_model "github.com/oj-lab/platform/models/judge"
	problem_model "github.com/oj-lab/platform/models/problem"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
)

func GetProblemInfoList(
	ctx context.Context,
	account string,
	options problem_model.GetProblemOptions,
) ([]problem_model.ProblemInfoView, int64, error) {
	options.Selection = problem_model.ProblemInfoSelection

	db := gorm_agent.GetDefaultDB()
	total, err := problem_model.CountProblemByOptions(db, options)
	if err != nil {
		return nil, 0, err
	}
	problemList, err := problem_model.GetProblemListByOptions(db, options)
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
				problemList[i].Solved = func() *bool { b := true; return &b }()
			} else {
				problemList[i].Solved = func() *bool { b := false; return &b }()
			}
		}
	}

	problemInfoList := problem_model.GetProblemInfoViewList(problemList)

	return problemInfoList, total, nil
}
