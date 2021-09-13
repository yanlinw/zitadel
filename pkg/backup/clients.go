package backup

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getS3Client(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	region string,
) *s3.Client {

	reg := region
	if region == "" {
		reg = "us-east-1"
	}

	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint, // or where ever you ran minio
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:           reg,
		Credentials:      credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		EndpointResolver: staticResolver,
	}

	return s3.NewFromConfig(cfg)
}
