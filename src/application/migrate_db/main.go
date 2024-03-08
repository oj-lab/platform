package main

import (
	"github.com/OJ-lab/oj-lab-services/src/core"
	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
)

func main() {
	db := gormAgent.GetDefaultDB()
	err := db.AutoMigrate(&model.User{}, &model.Problem{}, &model.JudgeTaskSubmission{}, &model.Judger{})
	if err != nil {
		panic("failed to migrate database")
	}

	// data init migrate to test

	// description := `Write a program that prints "Hello World!".`
	// mapper.CreateProblem(db, model.Problem{
	// 	Slug:        "hello-world",
	// 	Title:       "Hello World!",
	// 	Description: &description,
	// 	Tags: []*model.AlgorithmTag{
	// 		{Name: "Primer"},
	// 	},
	// })

	// description = `Calculate A + B, print the result.`
	// mapper.CreateProblem(db, model.Problem{
	// 	Slug:        "a-plus-b",
	// 	Title:       "A + B",
	// 	Description: &description,
	// 	Tags: []*model.AlgorithmTag{
	// 		{Name: "Primer"},
	// 		{Name: "Math"},
	// 	},
	// })

	// mapper.CreateUser(db, model.User{
	// 	Name:     "admin",
	// 	Account:  "admin",
	// 	Password: func() *string { s := "admin"; return &s }(),
	// 	Roles: []*model.Role{
	// 		{Name: "admin"},
	// 	},
	// })

	// mapper.CreateUser(db, model.User{
	// 	Name:     "anonymous",
	// 	Account:  "anonymous",
	// 	Password: func() *string { s := "anonymous"; return &s }(),
	// 	Roles: []*model.Role{
	// 		{Name: "anonymous"},
	// 	},
	// })

	core.AppLogger().Info("migrate tables success")
}
