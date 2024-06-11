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
	problem_model "github.com/oj-lab/oj-lab-platform/models/problem"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
	minio_agent "github.com/oj-lab/oj-lab-platform/modules/agent/minio"
	"github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/oj-lab/oj-lab-platform/modules/log"
	yaml "gopkg.in/yaml.v2"
)

var ctx = context.Background()

func main() {
	db := gorm_agent.GetDefaultDB()
	minioClient := minio_agent.GetMinioClient()

	// Read package files
	// Search Problem under packagePath
	// 1. parse problem path as `slug`,
	//    parse problem.yaml's name as `title`,
	//    parse problem.md as description.
	// 2. insert object into minio storage.
	var (
		packagePath string = path.Join(config.Workdir, "problem_packages")
		title       string
		slug        string
	)
	err := filepath.Walk(packagePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.AppLogger().WithError(err).Error("Walk package path failed")
			return err
		}
		if info == nil {
			return fmt.Errorf("file info is nil")
		}
		if info.IsDir() {
			return nil
		}
		relativePath := strings.Replace(path, packagePath, "", 1)
		log.AppLogger().WithField("relativePath", relativePath).Debug("Read file from package")
		if filepath.Base(relativePath) == "problem.yaml" {
			resultMap := make(map[string]interface{})
			yamlFile, err := os.ReadFile(path)
			if err != nil {
				log.AppLogger().WithError(err).Error("Read problem.yaml failed")
			}
			err = yaml.Unmarshal(yamlFile, &resultMap)
			if err != nil {
				log.AppLogger().WithError(err).Error("Unmarshal problem.yaml failed")
			}
			title = resultMap["name"].(string)
			if title == "" {
				log.AppLogger().Error("Problem title is empty")
			}
			slug = strings.Split(relativePath, "/")[1]
			log.AppLogger().WithField("title", title).WithField("slug", slug).Debug("Read problem.yaml")
		}
		if filepath.Base(relativePath) == "problem.md" {
			content, err := os.ReadFile(path)
			if err != nil {
				log.AppLogger().WithError(err).Error("Read problem.md failed")
			}
			description := string(content)
			log.AppLogger().WithField("description", description).Debug("Read problem.md")
			err = problem_model.CreateProblem(db, problem_model.Problem{
				Slug:        slug,
				Title:       title,
				Description: &description,
				Tags: []*problem_model.AlgorithmTag{
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
			log.AppLogger().WithError(err).Error("Put object to minio failed")
		}
		return err
	})
	if err != nil {
		panic(err)
	}

	log.AppLogger().Info("Problem loaded")
}
