package backup

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"strings"
)

func ListS3AssetBuckets(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	prefix string,
) ([]string, error) {

	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey)

	input := &s3.ListBucketsInput{}

	result, err := s3Client.ListBuckets(context.Background(), input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	buckets := make([]string, 0)
	for _, bucket := range result.Buckets {
		name := *bucket.Name
		if strings.HasPrefix(name, prefix) {
			buckets = append(buckets, *bucket.Name)
		}
	}

	return buckets, nil
}

func DeleteS3BucketIfExists(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	bucket string,
) error {
	ctx := context.Background()
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey)

	listObjects, err := s3Client.ListObjectsV2(
		ctx,
		&s3.ListObjectsV2Input{
			Bucket: &bucket,
		},
	)
	if err != nil {
		return err
	}
	for _, object := range listObjects.Contents {
		if _, err := s3Client.DeleteObject(
			ctx,
			&s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    object.Key,
			}); err != nil {
			return err
		}
	}

	_, err = s3Client.DeleteBucket(context.Background(),
		&s3.DeleteBucketInput{
			Bucket: &bucket,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func CreateS3BucketIfNotExists(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	bucket string,
) error {
	ctx := context.Background()
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey)
	input := &s3.CreateBucketInput{
		Bucket: &bucket,
	}

	currentBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	found := false
	for _, currentBucket := range currentBuckets.Buckets {
		if bucket == *currentBucket.Name {
			found = true
		}
	}

	if !found {
		_, err := s3Client.CreateBucket(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}

func CleanS3BucketsWithPrefix(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	prefix string,
) error {
	buckets, err := ListS3AssetBuckets(endpoint, accessKeyID, secretAccessKey, prefix)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if err := DeleteS3BucketIfExists(endpoint, accessKeyID, secretAccessKey, bucket); err != nil {
			return err
		}
	}
	return nil
}

func EnsureS3BucketsWithPrefix(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	buckets []string,
	prefix string,
) error {
	currentBuckets, err := ListS3AssetBuckets(endpoint, accessKeyID, secretAccessKey, prefix)
	if err != nil {
		return err
	}

	createBuckets := make([]string, 0)
	deleteBuckets := make([]string, 0)

	for _, currentBucket := range currentBuckets {
		necessary := false
		for _, bucket := range buckets {
			if currentBucket == bucket {
				necessary = true
			}
		}
		if !necessary {
			deleteBuckets = append(deleteBuckets, currentBucket)
		}
	}
	for _, bucket := range buckets {
		found := false
		for _, currentBucket := range currentBuckets {
			if bucket == currentBucket {
				found = true
			}
		}
		if !found {
			createBuckets = append(createBuckets, bucket)
		}
	}

	for _, createBucket := range createBuckets {
		if err := CreateS3BucketIfNotExists(endpoint, accessKeyID, secretAccessKey, createBucket); err != nil {
			return err
		}
	}

	for _, deleteBucket := range deleteBuckets {
		if err := DeleteS3BucketIfExists(endpoint, accessKeyID, secretAccessKey, deleteBucket); err != nil {
			return err
		}
	}

	return nil
}
