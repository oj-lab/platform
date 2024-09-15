package problem_model

import (
	"encoding/json"
	"fmt"
	"testing"

	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
)

func TestProblemDB(t *testing.T) {
	db := gorm_agent.GetDefaultDB()
	description := "Given two integer A and B, please output the answer of A+B."
	problem := Problem{
		Slug:        "a-plus-b-problem",
		Title:       "A+B Problem",
		Description: &description,
		Tags:        []*ProblemTag{{Name: "tag1"}, {Name: "tag2"}},
	}

	err := CreateProblem(db, problem)
	if err != nil {
		t.Error(err)
	}

	dbProblem, err := GetProblem(db, problem.Slug)
	if err != nil {
		t.Error(err)
	}

	problemJson, err := json.MarshalIndent(dbProblem, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemJson))

	problemOption := GetProblemOptions{
		Selection: ProblemInfoSelection,
		Tags:      []*ProblemTag{{Name: "tag1"}},
		Slug:      problem.Slug,
	}

	problemList, err := GetProblemListByOptions(db, problemOption)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("problemList: %v\n", problemList)
	if len(problemList) != 1 {
		t.Error("problemCount should be 1")
	}

	problemListJson, err := json.MarshalIndent(problemList, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemListJson))

	err = DeleteProblem(db, problem.Slug)
	if err != nil {
		t.Error(err)
	}
}
