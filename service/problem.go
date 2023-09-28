package service

import (
	"context"

	"github.com/OJ-lab/oj-lab-services/core/agent/judger"
	"github.com/OJ-lab/oj-lab-services/service/business"
	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func GetProblem(ctx context.Context, slug string) (*model.Problem, error) {
	problem, err := mapper.GetProblem(slug)
	if err != nil {
		return nil, err
	}
	return problem, nil
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

func Judge(ctx context.Context, slug string, judgeRequest judger.JudgeRequest) (
	[]map[string]interface{}, error,
) {
	body, err := judger.PostJudgeSync(slug, judgeRequest)
	if err != nil {
		return nil, err
	}

	return body, nil
}
