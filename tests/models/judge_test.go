package models_test

import (
	"testing"

	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func TestJudgeDB(t *testing.T) {
	db := gorm_agent.GetDefaultDB()
	problem := &problem_model.Problem{
		Slug: "test-judge-db-problem",
	}
	var err error
	err = problem_model.CreateProblem(db, *problem)
	if err != nil {
		t.Error(err)
	}
	judge := &judge_model.Judge{
		Language:    judge_model.ProgrammingLanguageCpp,
		ProblemSlug: problem.Slug,
	}
	judge, err = judge_model.CreateJudge(db, *judge)
	if err != nil {
		t.Error(err)
	}

	judgeResult := &judge_model.JudgeResult{
		JudgeUID: judge.UID,
		Verdict:  judge_model.JudgeVerdictAccepted,
	}
	_, err = judge_model.CreateJudgeResult(db, *judgeResult)
	if err != nil {
		t.Error(err)
	}

	judge, err = judge_model.GetJudge(db, judge.UID)
	if err != nil {
		t.Error(err)
	}
	t.Log(judge)
}
