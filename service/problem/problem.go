package problem

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	minioAgent "github.com/OJ-lab/oj-lab-services/packages/agent/minio"
	"github.com/OJ-lab/oj-lab-services/packages/mapper"
	"github.com/gin-gonic/gin"
)

func GetProblemInfo(ctx *gin.Context) {
	slug := ctx.Param("slug")
	problem, err := mapper.GetProblem(slug)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{
		"slug":        problem.Slug,
		"title":       problem.Title,
		"description": problem.Description,
		"tags":        mapper.GetTagsList(problem),
	})
}

func PutProblemPackage(ctx *gin.Context) {
	slug := ctx.Param("slug")
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}
	zipFile := "/tmp/" + slug + ".zip"
	if err := ctx.SaveUploadedFile(file, zipFile); err != nil {
		ctx.Error(err)
		return
	}

	// unzip package
	targetDir := "/tmp/" + slug
	err = os.RemoveAll(targetDir)
	if err != nil {
		ctx.Error(err)
		return
	}
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		ctx.Error(err)
		return
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			ctx.Error(err)
			return
		}
		defer rc.Close()

		path := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				ctx.Error(err)
				return
			}
			outFile, err := os.Create(path)
			if err != nil {
				ctx.Error(err)
				return
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, rc)
			if err != nil {
				ctx.Error(err)
				return
			}
		}
	}

	// put package to minio
	err = minioAgent.PutProblemPackage(slug, targetDir)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Done()
}
