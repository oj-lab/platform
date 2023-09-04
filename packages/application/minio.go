package application

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	minioEndpointProp        = "minio.endpoint"
	minioAccessKeyProp       = "minio.accessKeyID"
	minioSecretAccessKeyProp = "minio.secretAccessKey"
	minioUseSSLProp          = "minio.useSSL"
	minioRegionProp          = "minio.region"
)

var (
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	region          string
	minioClient     *minio.Client
)

func init() {
	endpoint = AppConfig.GetString(minioEndpointProp)
	accessKeyID = AppConfig.GetString(minioAccessKeyProp)
	secretAccessKey = AppConfig.GetString(minioSecretAccessKeyProp)
	useSSL = AppConfig.GetBool(minioUseSSLProp)
	region = AppConfig.GetString(minioRegionProp)
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
			panic("failed to connect minio")
		}
	}

	return minioClient
}
