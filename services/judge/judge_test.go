package judge_service

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/oj-lab/oj-lab-platform/models"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func TestGetJudge(t *testing.T) {
	db := gorm_agent.GetDefaultDB()
	problem := &problem_model.Problem{
		Slug: "test-judge-service",
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

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	response_judge, err := GetJudge(ctx, judge.UID)
	if err != nil {
		t.Error(err)
	}
	asserts := assert.New(t)
	asserts.Equal(judge.ProblemSlug, response_judge.ProblemSlug)
	asserts.Equal(judge.Language, response_judge.Language)
}

func TestCreateJudge(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	judge := &judge_model.Judge{
		Language:    judge_model.ProgrammingLanguageCpp,
		ProblemSlug: "test-judge-service",
	}
	_, err := CreateJudge(ctx, *judge)
	if err != nil {
		t.Error(err)
	}

	db := gorm_agent.GetDefaultDB()
	judges, err := judge_model.GetJudgeListByOptions(db,
		judge_model.GetJudgeOptions{OrderByColumns: []models.OrderByColumnOption{{Column: "create_at", Desc: true}}})
	if err != nil || len(judges) == 0 {
		t.Error(err)
	}
	insert_judge, err := judge_model.GetJudge(db, judges[0].UID)
	if err != nil {
		t.Error(err)
	}
	asserts := assert.New(t)
	asserts.Equal(judge.ProblemSlug, insert_judge.ProblemSlug)
	asserts.Equal(judge.Language, insert_judge.Language)
}

func TestUpsertJudgeCache(t *testing.T) {
	// previous WA || later WA ||  previous AC || later AC || GetBeforeSubmission
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	db := gorm_agent.GetDefaultDB()

	user := user_model.User{
		Account:  "test-upserJudgeCache-user",
		Password: func() *string { s := ""; return &s }(),
	}
	problem := problem_model.Problem{
		Slug: "test-upserJudgeCache-problem",
	}

	baseACJudge := &judge_model.Judge{
		UserAccount: user.Account,
		User:        user,
		ProblemSlug: problem.Slug,
		Problem:     problem,
		Verdict:     judge_model.JudgeVerdictAccepted,
		Status:      judge_model.JudgeStatusFinished,
	}
	baseACJudge.UID = uuid.New()
	baseACJudge, err := CreateJudge(ctx, *baseACJudge)
	if err != nil {
		t.Error(err)
	}
	time1 := time.Unix(int64(1000), 0)
	baseACJudge.MetaFields.CreateAt = &time1
	err = judge_model.UpdateJudge(db, *baseACJudge)
	if err != nil {
		t.Error(err)
	}
	err = UpsertJudgeCache(ctx, baseACJudge.UID, judge_model.JudgeVerdictAccepted)
	if err != nil {
		t.Error(err)
	}

	preWAJudge := &judge_model.Judge{
		UserAccount: user.Account,
		User:        user,
		ProblemSlug: problem.Slug,
		Problem:     problem,
		Verdict:     judge_model.JudgeVerdictWrongAnswer,
		Status:      judge_model.JudgeStatusFinished,
	}
	preWAJudge, err = CreateJudge(ctx, *preWAJudge)
	if err != nil {
		t.Error(err)
	}
	time2 := time.Unix(int64(998), 0)
	preWAJudge.MetaFields.CreateAt = &time2
	err = judge_model.UpdateJudge(db, *preWAJudge)
	if err != nil {
		t.Error(err)
	}
	err = UpsertJudgeCache(ctx, preWAJudge.UID, judge_model.JudgeVerdictWrongAnswer)
	if err != nil {
		t.Error(err)
	}

	laterWAJudge := preWAJudge
	laterWAJudge, err = CreateJudge(ctx, *laterWAJudge)
	if err != nil {
		t.Error(err)
	}
	time3 := time.Unix(int64(1001), 0)
	laterWAJudge.MetaFields.CreateAt = &time3
	err = judge_model.UpdateJudge(db, *laterWAJudge)
	if err != nil {
		t.Error(err)
	}
	err = UpsertJudgeCache(ctx, laterWAJudge.UID, judge_model.JudgeVerdictWrongAnswer)
	if err != nil {
		t.Error(err)
	}

	rankCache, err := judge_model.GetJudgeRankCache(db, user.Account)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rankCache)
	scoreCacheCache, err := judge_model.GetJudgeScoreCache(db, user.Account, problem.Slug)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(scoreCacheCache)

	asserts := assert.New(t)
	asserts.Equal(rankCache.Points, 1)
	asserts.Equal(rankCache.TotalSubmissions, 2)
	asserts.Equal(scoreCacheCache.SubmissionCount, 2)
	asserts.Equal(scoreCacheCache.SolveTime, baseACJudge.CreateAt)

	preACJudge := baseACJudge
	preACJudge, err = CreateJudge(ctx, *preACJudge)
	if err != nil {
		t.Error(err)
	}
	time4 := time.Unix(int64(999), 0)
	preACJudge.MetaFields.CreateAt = &time4
	err = judge_model.UpdateJudge(db, *preACJudge)
	if err != nil {
		t.Error(err)
	}
	err = UpsertJudgeCache(ctx, preACJudge.UID, judge_model.JudgeVerdictAccepted)
	if err != nil {
		t.Error(err)
	}
	laterACJudge := baseACJudge
	laterACJudge, err = CreateJudge(ctx, *laterACJudge)
	if err != nil {
		t.Error(err)
	}
	time5 := time.Unix(int64(1002), 0)
	laterACJudge.MetaFields.CreateAt = &time5
	err = judge_model.UpdateJudge(db, *laterACJudge)
	if err != nil {
		t.Error(err)
	}
	err = UpsertJudgeCache(ctx, laterACJudge.UID, judge_model.JudgeVerdictAccepted)
	if err != nil {
		t.Error(err)
	}

	rankCache, err = judge_model.GetJudgeRankCache(db, user.Account)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rankCache)
	scoreCacheCache, err = judge_model.GetJudgeScoreCache(db, user.Account, problem.Slug)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(scoreCacheCache)
	asserts.Equal(rankCache.Points, 1)
	asserts.Equal(rankCache.TotalSubmissions, 2)
	asserts.Equal(scoreCacheCache.SubmissionCount, 2)
	asserts.Equal(scoreCacheCache.SolveTime, preACJudge.CreateAt)

	submissionCount, err := judge_model.GetBeforeSubmission(db, *preACJudge)
	if err != nil {
		t.Error(err)
	}
	asserts.Equal(submissionCount, 2)
	submissionCount, err = judge_model.GetBeforeSubmission(db, *laterACJudge)
	if err != nil {
		t.Error(err)
	}
	asserts.Equal(submissionCount, 5)
}
