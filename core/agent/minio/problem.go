package minio

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

func PutProblemPackage(slug string, pkgDir string) error {
	ctx := context.Background()
	minioClient := GetMinioClient()

	// remove old package
	objectsCh := minioClient.ListObjects(ctx, GetBucketName(), minio.ListObjectsOptions{
		Prefix:    slug,
		Recursive: true,
	})
	for objInfo := range objectsCh {
		if objInfo.Err != nil {
			return objInfo.Err
		}

		err := minioClient.RemoveObject(ctx, GetBucketName(), objInfo.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}

	filepath.Walk(pkgDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relativePath := filepath.Join(slug, strings.Replace(path, pkgDir, "", 1))

		_, minioErr := minioClient.FPutObject(ctx, GetBucketName(),
			relativePath,
			path,
			minio.PutObjectOptions{})
		if minioErr != nil {
			log.Fatalln(minioErr)
		}
		return minioErr
	})

	return nil
}
