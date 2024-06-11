package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func TestProblemDB(t *testing.T) {
	db := gorm_agent.GetDefaultDB()
	description := "Given two integer A and B, please output the answer of A+B."
	problem := problem_model.Problem{
		Slug:        "a-plus-b-problem",
		Title:       "A+B Problem",
		Description: &description,
		Tags:        []*problem_model.AlgorithmTag{{Name: "tag1"}, {Name: "tag2"}},
	}

	err := problem_model.CreateProblem(db, problem)
	if err != nil {
		t.Error(err)
	}

	dbProblem, err := problem_model.GetProblem(db, problem.Slug)
	if err != nil {
		t.Error(err)
	}

	problemJson, err := json.MarshalIndent(dbProblem, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemJson))

	problemOption := problem_model.GetProblemOptions{
		Selection: problem_model.ProblemInfoSelection,
		Tags:      []*problem_model.AlgorithmTag{{Name: "tag1"}},
		Slug:      &problem.Slug,
	}

	problemList, problemCount, err := problem_model.GetProblemListByOptions(db, problemOption)
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

	err = problem_model.DeleteProblem(db, problem.Slug)
	if err != nil {
		t.Error(err)
	}
}
