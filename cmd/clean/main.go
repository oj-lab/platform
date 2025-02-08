package main

import (
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"
	judge_model "github.com/oj-lab/platform/models/judge"
	problem_model "github.com/oj-lab/platform/models/problem"
	user_model "github.com/oj-lab/platform/models/user"
	casbin_agent "github.com/oj-lab/platform/modules/agent/casbin"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	minio_agent "github.com/oj-lab/platform/modules/agent/minio"
	redis_agent "github.com/oj-lab/platform/modules/agent/redis"
)

func clearCasbin() {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	enforcer.ClearPolicy() // no err return
	slog.Info("Clear Casbin success")
}

func removeMinioObjects() {
	minioClient := minio_agent.GetMinioClient()
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		opts := minio.ListObjectsOptions{Recursive: true}
		for object := range minioClient.ListObjects(context.Background(), minio_agent.GetBucketName(), opts) {
			if object.Err != nil {
				slog.With("err", object.Err).Error("Get object error")
			}
			objectsCh <- object
		}
	}()

	errorCh := minioClient.RemoveObjects(context.Background(), minio_agent.GetBucketName(), objectsCh, minio.RemoveObjectsOptions{})
	for e := range errorCh {
		slog.With("err", e.Err).Error("Failed to remove " + e.ObjectName)
	}

	slog.Info("Remove Minio Objects success")
}

func clearRedis() {
	ctx := context.Background()
	redis_agent := redis_agent.GetDefaultRedisClient()
	err := redis_agent.FlushDB(ctx).Err()
	if err != nil {
		slog.With("err", err).Error("Failed to clear redis")
	}

	slog.Info("Clear Redis success")
}
func clearDB() {
	db := gorm_agent.GetDefaultDB()

	err := db.Migrator().DropTable(
		&user_model.User{},
		&problem_model.Problem{},
		&problem_model.ProblemTag{},
		&judge_model.Judge{},
		&judge_model.JudgeResult{},
		&judge_model.JudgeScoreCache{},
		&judge_model.JudgeRankCache{},
		"problem_problem_tags",
		"casbin_rule",
	)

	if err != nil {
		panic("failed to drop tables")
	}

	slog.Info("Clear DB success")
}

func main() {
	removeMinioObjects()
	clearCasbin()
	clearRedis()
	clearDB()

	println("data clean success")
}
