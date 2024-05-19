package minioAgent

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

func PutLocalObjects(ctx context.Context, rootName, localPath string) error {
	minioClient := GetMinioClient()
	bucketName := GetBucketName()

	// Remove old
	objectsCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    rootName,
		Recursive: true,
	})
	for objInfo := range objectsCh {
		if objInfo.Err != nil {
			return objInfo.Err
		}

		err := minioClient.RemoveObject(ctx, bucketName, objInfo.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}

	err := filepath.Walk(localPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		objectName := filepath.Join(rootName, strings.Replace(path, localPath, "", 1))
		_, err = minioClient.FPutObject(ctx, bucketName,
			objectName, path,
			minio.PutObjectOptions{})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
