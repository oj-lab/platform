package mapper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/packages/model"
)

func TestProblemMapper(t *testing.T) {
	problem := model.Problem{
		Slug:        "a+b-problem",
		Title:       "A+B Problem",
		Description: "Given two integer A and B, please output the answer of A+B.",
		Tags:        []*model.AlgorithmTag{{Slug: "tag1"}, {Slug: "tag2"}},
	}
	err := CreateProblem(problem)
	if err != nil {
		t.Error(err)
	}

	dbProblem, err := GetProblem(problem.Slug)
	if err != nil {
		t.Error(err)
	}

	problemJson, err := json.MarshalIndent(dbProblem, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemJson))

	problemOption := GetProblemOptions{
		Slug:  "",
		Title: "",
		Tags:  []*model.AlgorithmTag{{Slug: "tag1"}},
	}

	dbProblems, problemCount, err := GetProblemByOptions(problemOption)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("problemCount: %d\n", problemCount)
	if problemCount != 1 {
		t.Error("problemCount should be 1")
	}

	problemsJson, err := json.MarshalIndent(dbProblems, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemsJson))
}
