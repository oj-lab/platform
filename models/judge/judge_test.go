package judge_model

import (
	"encoding/json"
	"testing"

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
	judge := &Judge{
		Language:    ProgrammingLanguageCpp,
		ProblemSlug: problem.Slug,
	}
	judge, err = CreateJudge(db, *judge)
	if err != nil {
		t.Error(err)
	}

	judgeResult := &JudgeResult{
		JudgeUID: judge.UID,
		Verdict:  JudgeVerdictAccepted,
	}
	_, err = CreateJudgeResult(db, *judgeResult)
	if err != nil {
		t.Error(err)
	}

	judge, err = GetJudge(db, judge.UID)
	if err != nil {
		t.Error(err)
	}
	judgeJson, err := json.MarshalIndent(judge, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v\n", string(judgeJson))
}
