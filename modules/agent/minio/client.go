package minio_agent

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/oj-lab/oj-lab-platform/modules/config"
	"github.com/oj-lab/oj-lab-platform/modules/log"
)

const (
	minioEndpointProp        = "minio.endpoint"
	minioAccessKeyProp       = "minio.accessKeyID"
	minioSecretAccessKeyProp = "minio.secretAccessKey"
	minioUseSSLProp          = "minio.useSSL"
	minioRegionProp          = "minio.region"
	minioBucketNameProp      = "minio.bucketName"
)

var (
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	region          string
	minioClient     *minio.Client
	bucketName      string
)

func init() {
	endpoint = config.AppConfig.GetString(minioEndpointProp)
	accessKeyID = config.AppConfig.GetString(minioAccessKeyProp)
	secretAccessKey = config.AppConfig.GetString(minioSecretAccessKeyProp)
	useSSL = config.AppConfig.GetBool(minioUseSSLProp)
	region = config.AppConfig.GetString(minioRegionProp)
	bucketName = config.AppConfig.GetString(minioBucketNameProp)
}

func GetBucketName() string {
	return bucketName
}

func GetMinioClient() *minio.Client {
	if minioClient == nil {
		var err error
		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
			Region: region,
		})
		if err != nil {
			panic("failed to connect minio client")
		}
		ctx := context.Background()

		exists, err := minioClient.BucketExists(ctx, bucketName)
		if err == nil && exists {
			log.AppLogger().WithField("bucket", bucketName).Info("Bucket already exists")
			return minioClient
		}

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.AppLogger().WithError(err).
				WithField("bucket", bucketName).Error("Failed to create bucket")
		} else {
			log.AppLogger().WithField("bucket", bucketName).Info("Successfully created bucket")
		}
	}

	return minioClient
}
