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

	packagePath := path.Join(config_module.ProjectRoot(), "problem-packages/icpc")

	// Load Dirs under `packagePath`
	problemPackageDirs, err := os.ReadDir(packagePath)
	if err != nil {
		log_module.AppLogger().WithError(err).Error("Read package path failed")
		panic(err)
	}
	for _, problemPackageDir := range problemPackageDirs {
		if !problemPackageDir.IsDir() {
			continue
		}

		var (
			title               string
			slug                string
			testCaseID          = 1
			description         string
			testCaseDescription = "## Examples\n"
		)

		problemPackagePath := path.Join(packagePath, problemPackageDir.Name())
		err := filepath.Walk(problemPackagePath, func(path string, info fs.FileInfo, err error) error {
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
			log_module.AppLogger().
				WithField("relativePath", relativePath).
				WithField("Ext", filepath.Ext(relativePath)).
				WithField("Dir", filepath.Dir(relativePath)).
				Debug("Read file from package")

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
				description = string(content)
				log_module.AppLogger().WithField("description", description).Debug("Read problem.md")
			}

			if filepath.Ext(relativePath) == ".in" && strings.HasSuffix(filepath.Dir(relativePath), "sample") {
				ansPath := strings.Replace(path, ".in", ".ans", 1)
				if _, err := os.Stat(ansPath); err != nil {
					log_module.AppLogger().WithField("path", ansPath).Error("Answer file not found")
					return nil
				}
				input, err := os.ReadFile(path)
				if err != nil {
					log_module.AppLogger().WithError(err).Error("Read input file failed")
					return nil
				}
				inputStr := strings.Trim(string(input), "\n")
				output, err := os.ReadFile(ansPath)
				if err != nil {
					log_module.AppLogger().WithError(err).Error("Read output file failed")
					return nil
				}
				outputStr := strings.Trim(string(output), "\n")
				testCaseDescription += fmt.Sprintf("\n### Example %d\n", testCaseID)
				testCaseDescription += fmt.Sprintf("Input\n```text\n%s\n```\n", string(inputStr))
				testCaseDescription += fmt.Sprintf("Output\n```text\n%s\n```\n", string(outputStr))
				testCaseID++
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
		description += testCaseDescription
		err = problem_model.CreateProblem(db, problem_model.Problem{
			Slug:        slug,
			Title:       title,
			Description: &description,
		})
		if err != nil {
			panic(err)
		}
	}

	log_module.AppLogger().Info("Problem loaded")
}
