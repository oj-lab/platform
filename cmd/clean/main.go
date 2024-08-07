package main

import (
	"context"

	"github.com/minio/minio-go/v7"
	casbin_agent "github.com/oj-lab/oj-lab-platform/modules/agent/casbin"

	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	user_model "github.com/oj-lab/oj-lab-platform/models/user"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	minio_agent "github.com/oj-lab/oj-lab-platform/modules/agent/minio"

	log_module "github.com/oj-lab/oj-lab-platform/modules/log"
)

func clearCasbin() {
	enforcer := casbin_agent.GetDefaultCasbinEnforcer()
	enforcer.ClearPolicy() // no err return
	log_module.AppLogger().Info("Clear Casbin success")
}

func removeMinioObjects() {
	minioClient := minio_agent.GetMinioClient()
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		opts := minio.ListObjectsOptions{Recursive: true}
		for object := range minioClient.ListObjects(context.Background(), minio_agent.GetBucketName(), opts) {
			if object.Err != nil {
				log_module.AppLogger().WithError(object.Err).Error("Get object error")
			}
			objectsCh <- object
		}
	}()

	errorCh := minioClient.RemoveObjects(context.Background(), minio_agent.GetBucketName(), objectsCh, minio.RemoveObjectsOptions{})
	for e := range errorCh {
		log_module.AppLogger().WithError(e.Err).Error("Failed to remove " + e.ObjectName)
	}

	log_module.AppLogger().Info("Remove Minio Objects success")
}

func clearDB() {
	db := gorm_agent.GetDefaultDB()

	err := db.Migrator().DropTable(
		&user_model.User{},
		&problem_model.Problem{},
		&problem_model.AlgorithmTag{},
		&judge_model.Judge{},
		&judge_model.JudgeResult{},
		"problem_algorithm_tags",
		"casbin_rule",
	)

	if err != nil {
		panic("failed to drop tables")
	}

	log_module.AppLogger().Info("Clear db success")

}
func main() {
	removeMinioObjects()
	clearCasbin()
	clearDB()

	println("data clean success")
}
