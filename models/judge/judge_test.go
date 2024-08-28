package judge_model

import (
	"encoding/json"
	"testing"
	"time"

	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
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

func TestJudgeScoreCacheDB(t *testing.T) {
	db := gorm_agent.GetDefaultDB()

	problem := &problem_model.Problem{
		Slug: "test-judgeScoreCache-db-problem",
	}
	var err error
	err = problem_model.CreateProblem(db, *problem)
	if err != nil {
		t.Error(err)
	}

	user := &user_model.User{
		Account: "test-judgeScoreCache-db-user",
	}
	_, err = user_model.CreateUser(db, *user)
	if err != nil {
		t.Error(err)
	}

	judgeScoreCache := NewJudgeScoreCache(user.Account, problem.Slug)
	_, err = CreateJudgeScoreCache(db, judgeScoreCache)
	if err != nil {
		t.Error(err)
	}

	now := time.Now()
	judgeScoreCache.IsAccepted = true
	judgeScoreCache.SolveTime = &now
	judgeScoreCache.SubmissionCount = 10
	err = UpdateJudgeScoreCache(db, judgeScoreCache)
	if err != nil {
		t.Error(err)
	}

	newjudgeScoreCache, err := GetJudgeScoreCache(db, user.Account, problem.Slug)
	if err != nil {
		t.Error(err)
	}

	judgeScoreCacheJson, err := json.MarshalIndent(newjudgeScoreCache, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v\n", string(judgeScoreCacheJson))
}
