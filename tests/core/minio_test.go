package core_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	minioAgent "github.com/OJ-lab/oj-lab-services/src/core/agent/minio"
	"github.com/minio/minio-go/v7"
)

func TestMinio(T *testing.T) {
	minioClient := minioAgent.GetMinioClient()
	bucketName := minioAgent.GetBucketName()

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists != nil && !exists {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload package files
	packagePath := "../data/packages/icpc/hello_world"
	err = filepath.Walk(packagePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil {
			return fmt.Errorf("file info is nil")
		}
		if info.IsDir() {
			return nil
		}
		relativePath := filepath.Join(filepath.Base(packagePath), strings.Replace(path, packagePath, "", 1))
		println(relativePath)
		// A presigned URL for uploading objects with PutObject
		presignedURL, err := minioClient.PresignedPutObject(ctx, bucketName, relativePath, time.Minute*1)
		if err != nil {
			return err
		}

		uploadingFile, err := os.Open(path)
		if err != nil {
			return err
		}
		fileInfo, err := uploadingFile.Stat()
		if err != nil {
			return err
		}
		// Read the file data
		fileData := make([]byte, fileInfo.Size())
		_, err = uploadingFile.Read(fileData)
		if err != nil {
			return err
		}

		// Upload the file by presigned URL
		httpClient := &http.Client{}
		req, err := http.NewRequest("PUT", presignedURL.String(), bytes.NewBuffer(fileData))

		if err != nil {
			return err
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to upload file: %s", resp.Status)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}
