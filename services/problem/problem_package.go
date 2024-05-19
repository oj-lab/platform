package problem

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"

	minioAgent "github.com/oj-lab/oj-lab-platform/modules/agent/minio"
)

func unzipProblemPackage(ctx context.Context, zipFile, targetDir string) error {
	err := os.RemoveAll(targetDir)
	if err != nil {
		return err
	}
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func putProblemPackage(ctx context.Context, slug string, pkgDir string) error {
	err := minioAgent.PutLocalObjects(ctx, slug, pkgDir)
	if err != nil {
		return err
	}

	return nil
}
