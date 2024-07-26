package judge_service

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
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
