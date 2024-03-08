package core_test

import (
	"testing"

	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"

	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
)

func TestPG(T *testing.T) {
	db := gormAgent.GetDefaultDB()

	description := `Write a program that prints "Hello World!".` // data migrate to test
	mapper.CreateProblem(db, model.Problem{
		Slug:        "hello-world",
		Title:       "Hello World!",
		Description: &description,
		Tags: []*model.AlgorithmTag{
			{Name: "Primer"},
		},
	})

	description = `Calculate A + B, print the result.`
	mapper.CreateProblem(db, model.Problem{
		Slug:        "a-plus-b",
		Title:       "A + B",
		Description: &description,
		Tags: []*model.AlgorithmTag{
			{Name: "Primer"},
			{Name: "Math"},
		},
	})

	mapper.CreateUser(db, model.User{
		Name:     "admin",
		Account:  "admin",
		Password: func() *string { s := "admin"; return &s }(),
		Roles: []*model.Role{
			{Name: "admin"},
		},
	})

	mapper.CreateUser(db, model.User{
		Name:     "anonymous",
		Account:  "anonymous",
		Password: func() *string { s := "anonymous"; return &s }(),
		Roles: []*model.Role{
			{Name: "anonymous"},
		},
	})
}
