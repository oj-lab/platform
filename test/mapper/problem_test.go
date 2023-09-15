package mapper_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/package/mapper"
	"github.com/OJ-lab/oj-lab-services/package/model"
)

func TestProblemMapper(t *testing.T) {
	problem := model.Problem{
		Slug:        "a+b-problem",
		Title:       "A+B Problem",
		Description: "Given two integer A and B, please output the answer of A+B.",
		Tags:        []*model.AlgorithmTag{{Slug: "tag1"}, {Slug: "tag2"}},
	}
	err := mapper.CreateProblem(problem)
	if err != nil {
		t.Error(err)
	}

	dbProblem, err := mapper.GetProblem(problem.Slug)
	if err != nil {
		t.Error(err)
	}

	problemJson, err := json.MarshalIndent(dbProblem, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(problemJson))

	problemOption := mapper.GetProblemOptions{
		Slug:  "",
		Title: "",
		Tags:  []*model.AlgorithmTag{{Slug: "tag1"}},
	}

	dbProblems, problemCount, err := mapper.GetProblemByOptions(problemOption)
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
