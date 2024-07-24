package main

import (
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	log_module "github.com/oj-lab/oj-lab-platform/modules/log"
)

func initDB() {
	db := gorm_agent.GetDefaultDB()
	err := db.AutoMigrate(
		&user_model.User{},
		&problem_model.Problem{},
		&judge_model.Judge{},
		&judge_model.JudgeResult{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	err = user_model.CreateUser(db, user_model.User{
		Name:     "root",
		Account:  "root",
		Password: func() *string { s := ""; return &s }(),
	})
	if err != nil {
		panic("failed to create admin user")
	}

	err = user_model.CreateUser(db, user_model.User{
		Name:     "anonymous",
		Account:  "anonymous",
		Password: func() *string { s := ""; return &s }(),
	})
	if err != nil {
		panic("failed to create anonymous user")
	}
	log_module.AppLogger().Info("migrate tables ans users success")
}
