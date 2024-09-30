package judge_model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	problem_model "github.com/oj-lab/platform/models/problem"
	user_model "github.com/oj-lab/platform/models/user"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	"github.com/stretchr/testify/assert"
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

func TestJudgeRank(t *testing.T) {
	db := gorm_agent.GetDefaultDB()

	var names []string
	var users []user_model.User

	for i := 0; i < 4; i++ {
		names = append(names, "test"+strconv.Itoa(i))
		users = append(users, user_model.User{
			Account:  names[i],
			Password: func() *string { s := ""; return &s }(),
		})
	}

	_, err := CreateJudgeRankCache(db, JudgeRankCache{
		UserAccount:      names[0],
		User:             users[0],
		Points:           10,
		TotalSubmissions: 10,
	})
	if err != nil {
		t.Error(err)
	}

	_, err = CreateJudgeRankCache(db, JudgeRankCache{
		UserAccount:      names[1],
		User:             users[1],
		Points:           10,
		TotalSubmissions: 11,
	})
	if err != nil {
		t.Error(err)
	}

	_, err = CreateJudgeRankCache(db, JudgeRankCache{
		UserAccount:      names[2],
		User:             users[2],
		Points:           10,
		TotalSubmissions: 12,
	})
	if err != nil {
		t.Error(err)
	}

	_, err = CreateJudgeRankCache(db, JudgeRankCache{
		UserAccount:      names[3],
		User:             users[3],
		Points:           9,
		TotalSubmissions: 10,
	})
	if err != nil {
		t.Error(err)
	}

	err = UpdateJudgeRankCache(db, JudgeRankCache{
		UserAccount:      names[3],
		User:             users[3],
		Points:           9,
		TotalSubmissions: 8,
	})
	assert.Error(t, err)
	assert.Equal(t, "TotalSubmissions must >= Points", err.Error())
	rankCache, err := GetJudgeRankCache(db, names[3])
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("rankCache: %v\n", rankCache)

	rankOption := GetRankCacheOptions{
		Selection: RankCacheInfoSelection,
	}
	rankList, err := GetRankCacheListByOptions(db, rankOption)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("rankList: %v\n", rankList)
	if len(rankList) != 4 {
		t.Error("rankCount should be 4")
	}

	offset := 2
	rankOption = GetRankCacheOptions{
		Selection: RankCacheInfoSelection,
		Offset:    &offset,
	}
	rankList, err = GetRankCacheListByOptions(db, rankOption)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("rankList: %v\n", rankList)
	if len(rankList) != 2 {
		t.Error("rankCount should be 2")
	}

	rankListJson, err := json.MarshalIndent(rankList, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(rankListJson))

	for i := 0; i < 4; i++ {
		err = DeleteJudgeRankCache(db, names[i])
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 4; i++ {
		err = user_model.DeleteUser(db, users[i].Account)
		if err != nil {
			t.Error(err)
		}
	}
}
