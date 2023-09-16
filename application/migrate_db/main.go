package main

import (
	"github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"github.com/sirupsen/logrus"
)

func main() {
	db := gorm.GetDefaultDB()
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
		Account:  "admin",
		Password: func() *string { s := "admin"; return &s }(),
	})

	logrus.Info("migrate tables success")
}
