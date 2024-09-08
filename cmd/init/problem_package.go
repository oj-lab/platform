package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	problem_model "github.com/oj-lab/platform/models/problem"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
	minio_agent "github.com/oj-lab/platform/modules/agent/minio"
	config_module "github.com/oj-lab/platform/modules/config"
	log_module "github.com/oj-lab/platform/modules/log"
	"gopkg.in/yaml.v2"
)

func loadProblemPackages(ctx context.Context) {
	db := gorm_agent.GetDefaultDB()
	minioClient := minio_agent.GetMinioClient()

	// Read package files
	// Search Problem under packagePath
	// 1. parse problem path as `slug`,
	//    parse problem.yaml's name as `title`,
	//    parse problem.md as description.
	// 2. insert object into minio storage.
	var (
		packagePath string = path.Join(config_module.ProjectRoot(), "problem_packages/icpc")
		title       string
		slug        string
	)
	err := filepath.Walk(packagePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log_module.AppLogger().WithError(err).Error("Walk package path failed")
			return err
		}
		if info == nil {
			return fmt.Errorf("file info is nil")
		}
		if info.IsDir() {
			return nil
		}
		relativePath := strings.Replace(path, packagePath, "", 1)
		log_module.AppLogger().WithField("relativePath", relativePath).Debug("Read file from package")
		if filepath.Base(relativePath) == "problem.yaml" {
			resultMap := make(map[string]interface{})
			yamlFile, err := os.ReadFile(path)
			if err != nil {
				log_module.AppLogger().WithError(err).Error("Read problem.yaml failed")
			}
			err = yaml.Unmarshal(yamlFile, &resultMap)
			if err != nil {
				log_module.AppLogger().WithError(err).Error("Unmarshal problem.yaml failed")
			}
			if resultMap["name"] == nil {
				log_module.AppLogger().Error("Problem name is nil")
				return nil
			}
			title = resultMap["name"].(string)
			if title == "" {
				log_module.AppLogger().Error("Problem title is empty")
			}
			slug = strings.Split(relativePath, "/")[1]
			log_module.AppLogger().WithField("title", title).WithField("slug", slug).Debug("Read problem.yaml")
		}
		if filepath.Base(relativePath) == "problem.md" {
			content, err := os.ReadFile(path)
			if err != nil {
				log_module.AppLogger().WithError(err).Error("Read problem.md failed")
			}
			description := string(content)
			log_module.AppLogger().WithField("description", description).Debug("Read problem.md")
			err = problem_model.CreateProblem(db, problem_model.Problem{
				Slug:        slug,
				Title:       title,
				Description: &description,
				Tags: []*problem_model.ProblemTag{
					{Name: "to-be-add"},
				},
			})
			if err != nil {
				return err
			}
		}

		_, err = minioClient.FPutObject(ctx, minio_agent.GetBucketName(),
			relativePath,
			path,
			minio.PutObjectOptions{})
		if err != nil {
			log_module.AppLogger().WithError(err).Error("Put object to minio failed")
		}
		return err
	})
	if err != nil {
		panic(err)
	}

	log_module.AppLogger().Info("Problem loaded")
}
