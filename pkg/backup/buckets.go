package backup

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"strings"
)

type ServiceAccount struct {
	ProjectID string `json:"project_id"`
}

func ListGCSAssetBuckets(
	serviceAccountJSON string,
	prefix string,
) ([]string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
	if err != nil {
		return nil, err
	}

	sa := &ServiceAccount{}
	if err := json.Unmarshal([]byte(serviceAccountJSON), sa); err != nil {
		return nil, err
	}

	buckets := make([]string, 0)
	it := client.Buckets(ctx, sa.ProjectID)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(attrs.Name, prefix) {
			buckets = append(buckets, attrs.Name)
		}
	}

	return buckets, nil
}

func DeleteGCSBucketIfExists(
	serviceAccountJSON string,
	bucket string,
) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
	if err != nil {
		return err
	}

	sa := &ServiceAccount{}
	if err := json.Unmarshal([]byte(serviceAccountJSON), sa); err != nil {
		return err
	}

	found := false
	it := client.Buckets(ctx, sa.ProjectID)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if attrs.Name == bucket {
			found = true
			break
		}
	}
	if found {
		return client.Bucket(bucket).Delete(ctx)
	}
	return nil
}

func CreateGCSBucketIfNotExists(
	serviceAccountJSON string,
	bucket string,
	location string,
) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
	if err != nil {
		return err
	}

	sa := &ServiceAccount{}
	if err := json.Unmarshal([]byte(serviceAccountJSON), sa); err != nil {
		return err
	}

	found := false
	it := client.Buckets(ctx, sa.ProjectID)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if attrs.Name == bucket {
			found = true
			break
		}
	}
	if !found {
		return client.Bucket(bucket).Create(ctx, sa.ProjectID, &storage.BucketAttrs{
			Name:         bucket,
			StorageClass: "STANDARD",
			Location:     location,
		})
	}
	return nil
}

func CleanGCSBucketsWithPrefix(
	serviceAccountJSON string,
	prefix string,
) error {
	buckets, err := ListGCSAssetBuckets(serviceAccountJSON, prefix)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if err := DeleteGCSBucketIfExists(serviceAccountJSON, bucket); err != nil {
			return err
		}
	}
	return nil
}

func EnsureGCSBucketsWithPrefix(
	serviceAccountJSON string,
	buckets []string,
	prefix string,
	location string,
) error {
	currentBuckets, err := ListGCSAssetBuckets(serviceAccountJSON, prefix)
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
		if err := CreateGCSBucketIfNotExists(serviceAccountJSON, createBucket, location); err != nil {
			return err
		}
	}

	for _, deleteBucket := range deleteBuckets {
		if err := DeleteGCSBucketIfExists(serviceAccountJSON, deleteBucket); err != nil {
			return err
		}
	}

	return nil
}

func ListS3AssetBuckets(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	prefix string,
	region string,
) ([]string, error) {

	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey, region)

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
	region string,
) error {
	ctx := context.Background()
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey, region)

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
	region string,
) error {
	ctx := context.Background()
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey, region)
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
	region string,
) error {
	buckets, err := ListS3AssetBuckets(endpoint, accessKeyID, secretAccessKey, prefix, region)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if err := DeleteS3BucketIfExists(endpoint, accessKeyID, secretAccessKey, bucket, region); err != nil {
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
	region string,
) error {
	currentBuckets, err := ListS3AssetBuckets(endpoint, accessKeyID, secretAccessKey, prefix, region)
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
		if err := CreateS3BucketIfNotExists(endpoint, accessKeyID, secretAccessKey, createBucket, region); err != nil {
			return err
		}
	}

	for _, deleteBucket := range deleteBuckets {
		if err := DeleteS3BucketIfExists(endpoint, accessKeyID, secretAccessKey, deleteBucket, region); err != nil {
			return err
		}
	}

	return nil
}
