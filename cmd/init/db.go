package main

import (
	"fmt"

	judge_model "github.com/oj-lab/platform/models/judge"
	problem_model "github.com/oj-lab/platform/models/problem"
	user_model "github.com/oj-lab/platform/models/user"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	log_module "github.com/oj-lab/platform/modules/log"
)

func initDB() {
	db := gorm_agent.GetDefaultDB()
	err := db.AutoMigrate(
		&user_model.User{},
		&problem_model.Problem{},
		&judge_model.Judge{},
		&judge_model.JudgeResult{},
		&judge_model.JudgeScoreCache{},
		&judge_model.JudgeRankCache{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	_, err = user_model.CreateUser(db, user_model.User{
		Name:     "root",
		Account:  "root",
		Password: func() *string { s := ""; return &s }(),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create root user: %v", err))
	}

	_, err = user_model.CreateUser(db, user_model.User{
		Name:     "anonymous",
		Account:  "anonymous",
		Password: func() *string { s := ""; return &s }(),
	})

	if err != nil {
		panic(fmt.Sprintf("failed to create anonymous user: %v", err))
	}

	log_module.AppLogger().Info("migrate tables ans users success")
}
