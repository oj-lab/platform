package service

import (
	"github.com/OJ-lab/oj-lab-services/package/agent/judger"
	"github.com/OJ-lab/oj-lab-services/package/agent/minio"
	"github.com/OJ-lab/oj-lab-services/package/mapper"
	"github.com/OJ-lab/oj-lab-services/package/model"
	"github.com/OJ-lab/oj-lab-services/service/business"
)

func GetProblemInfo(slug string) (*model.Problem, error) {
	problem, err := mapper.GetProblem(slug)
	if err != nil {
		return nil, err
	}
	return problem, nil
}

func PutProblemPackage(slug, zipFile string) error {
	targetDir := "/tmp/" + slug
	err := business.UnzipProblemPackage(zipFile, targetDir)
	if err != nil {
		return err
	}

	err = minio.PutProblemPackage(slug, targetDir)
	if err != nil {
		return err
	}

	return nil
}

func Judge(slug string, judgeRequest judger.JudgeRequest) (
	[]map[string]interface{}, error,
) {
	body, err := judger.PostJudgeSync(slug, judgeRequest)
	if err != nil {
		return nil, err
	}

	return body, nil
}
