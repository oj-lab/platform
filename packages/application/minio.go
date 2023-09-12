package application

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	endpoint = AppConfig.GetString(minioEndpointProp)
	accessKeyID = AppConfig.GetString(minioAccessKeyProp)
	secretAccessKey = AppConfig.GetString(minioSecretAccessKeyProp)
	useSSL = AppConfig.GetBool(minioUseSSLProp)
	region = AppConfig.GetString(minioRegionProp)
	bucketName = AppConfig.GetString(minioBucketNameProp)
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

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", bucketName)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Successfully created %s\n", bucketName)
		}
	}

	return minioClient
}
