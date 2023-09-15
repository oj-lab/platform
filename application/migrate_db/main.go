package main

import (
	"github.com/OJ-lab/oj-lab-services/package/core"
	"github.com/OJ-lab/oj-lab-services/package/mapper"
	"github.com/OJ-lab/oj-lab-services/package/model"
	"github.com/sirupsen/logrus"
)

func main() {
	db := core.GetDefaultDB()
	err := db.AutoMigrate(&model.User{}, &model.Problem{})
	if err != nil {
		panic("failed to migrate database")
	}

	mapper.CreateProblem(model.Problem{
		Slug:        "hello-world",
		Title:       "Hello! { ... }",
		Description: `Write a program that prints "Hello! %s" to the standard output (stdout).`,
	})

	mapper.CreateUser(model.User{
		Name:     "admin",
		Password: func() *string { s := "admin"; return &s }(),
	})

	logrus.Info("migrate tables success")
}
