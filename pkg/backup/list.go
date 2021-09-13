package backup

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"strings"
)

func GCSListFilesWithFilter(serviceAccountJSON string, bucketName, name string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
	if err != nil {
		return nil, err
	}

	bkt := client.Bucket(bucketName)

	names := make([]string, 0)
	it := bkt.Objects(ctx, &storage.Query{Prefix: name + "/"})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		parts := strings.Split(attrs.Name, "/")
		found := false
		for _, name := range names {
			if len(parts) < 2 {
				continue
			}
			if name == parts[1] {
				found = true
			}
		}
		if !found {
			names = append(names, parts[1])
		}
	}

	return names, nil
}

func S3ListFilesWithFilter(accessKeyID, secretAccessKey string, endpoint string, bucketName, name string, region string) ([]string, error) {
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey, region)

	prefix := name + "/"
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}
	paginator := s3.NewListObjectsV2Paginator(s3Client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = 10
	})

	names := make([]string, 0)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, value := range output.Contents {
			if strings.HasPrefix(*value.Key, prefix) {
				parts := strings.Split(*value.Key, "/")
				if len(parts) < 2 {
					continue
				}
				found := false
				for _, name := range names {
					if name == parts[1] {
						found = true
					}
				}
				if !found {
					names = append(names, parts[1])
				}
				names = append(names)
			}

		}
	}
	return names, nil
}
