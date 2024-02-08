package service

import (
	"context"

	gormAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/src/service/business"
	"github.com/OJ-lab/oj-lab-services/src/service/mapper"
	"github.com/OJ-lab/oj-lab-services/src/service/model"
	"gorm.io/gorm"
)

func GetProblem(ctx context.Context, slug string) (*model.Problem, error) {
	db := gormAgent.GetDefaultDB()
	problem, err := mapper.GetProblem(db, slug)
	if err != nil {
		return nil, err
	}
	return problem, nil
}

func PutProblem(ctx context.Context, problem model.Problem) error {
	db := gormAgent.GetDefaultDB()
	err := mapper.CreateProblem(db, problem)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProblem(ctx context.Context, slug string) error {
	db := gormAgent.GetDefaultDB()
	err := mapper.DeleteProblem(db, slug)
	if err != nil {
		return err
	}
	return nil
}

func CheckProblemSlug(ctx context.Context, slug string) (bool, error) {
	db := gormAgent.GetDefaultDB()
	problem, err := mapper.GetProblem(db, slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}
	return problem == nil, nil
}

func GetProblemInfoList(ctx context.Context) ([]model.ProblemInfo, int64, error) {
	return business.GetProblemInfoList(ctx)
}

func PutProblemPackage(ctx context.Context, slug, zipFile string) error {
	localDir := "/tmp/" + slug
	err := business.UnzipProblemPackage(ctx, zipFile, localDir)
	if err != nil {
		return err
	}

	err = business.PutProblemPackage(ctx, slug, localDir)
	if err != nil {
		return err
	}

	return nil
}

// func Judge(ctx context.Context, slug string, code string, language string) (
// 	[]map[string]interface{}, error,
// ) {
// 	request := judger.JudgeRequest{
// 		Code:     code,
// 		Language: language,
// 	}
// 	responseBody, err := judger.PostJudgeSync(slug, request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return responseBody, nil
// }
