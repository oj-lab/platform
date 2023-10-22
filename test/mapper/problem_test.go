package mapper_test

import (
	"encoding/json"
	"fmt"
	"testing"

	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func TestProblemMapper(t *testing.T) {
	db := gormAgent.GetDefaultDB()
	description := "Given two integer A and B, please output the answer of A+B."
	problem := model.Problem{
		Slug:        "a-plus-b-problem",
		Title:       "A+B Problem",
		Description: &description,
		Tags:        []*model.AlgorithmTag{{Slug: "tag1"}, {Slug: "tag2"}},
	}

	err := mapper.CreateProblem(db, problem)
	if err != nil {
		t.Error(err)
	}

	dbProblem, err := mapper.GetProblem(db, problem.Slug)
	if err != nil {
		t.Error(err)
	}

	problemJson, err := json.MarshalIndent(dbProblem, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemJson))

	problemOption := mapper.GetProblemOptions{
		Selection: model.ProblemInfoSelection,
		Tags:      []*model.AlgorithmTag{{Slug: "tag1"}},
		Slug:      &problem.Slug,
	}

	problemList, problemCount, err := mapper.GetProblemListByOptions(db, problemOption)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("problemCount: %d\n", problemCount)
	if problemCount != 1 {
		t.Error("problemCount should be 1")
	}

	problemListJson, err := json.MarshalIndent(problemList, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemListJson))
}
