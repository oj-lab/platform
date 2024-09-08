package minio_agent

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	config_module "github.com/oj-lab/platform/modules/config"
	log_module "github.com/oj-lab/platform/modules/log"
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
	endpoint = config_module.AppConfig().GetString(minioEndpointProp)
	accessKeyID = config_module.AppConfig().GetString(minioAccessKeyProp)
	secretAccessKey = config_module.AppConfig().GetString(minioSecretAccessKeyProp)
	useSSL = config_module.AppConfig().GetBool(minioUseSSLProp)
	region = config_module.AppConfig().GetString(minioRegionProp)
	bucketName = config_module.AppConfig().GetString(minioBucketNameProp)
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
			log_module.AppLogger().WithField("bucket", bucketName).Info("Bucket already exists")
			return minioClient
		}

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log_module.AppLogger().WithError(err).
				WithField("bucket", bucketName).Error("Failed to create bucket")
		} else {
			log_module.AppLogger().WithField("bucket", bucketName).Info("Successfully created bucket")
		}
	}

	return minioClient
}
