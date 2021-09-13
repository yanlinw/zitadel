package backup

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ListS3Folders(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	bucketName string,
	path string,
	region string,
) ([]string, error) {
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey, region)
	ctx := context.Background()
	input := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &path,
	}

	output, err := s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0)
	for _, obj := range output.Contents {
		objects = append(objects, *obj.Key)
	}

	return objects, nil
}

func ListGCSFolders(
	saJSONPath string,
	bucketName string,
	path string,
) ([]string, error) {
	ctx := context.Background()
	data, err := ioutil.ReadFile(saJSONPath)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(data))
	if err != nil {
		return nil, err
	}

	iter := client.Bucket(bucketName).Objects(
		ctx,
		&storage.Query{
			Prefix:   path,
			Versions: false,
		},
	)
	objects := make([]string, 0)
	for {
		obj, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(obj.Name, path) {
			folder := strings.TrimPrefix(obj.Name, path+"/")
			for {
				folderT, _ := filepath.Split(folder)
				folder = strings.TrimSuffix(folderT, string(filepath.Separator))
				if !strings.Contains(folder, string(filepath.Separator)) {
					break
				}
			}

			alreadyListed := false
			for _, object := range objects {
				if object == folder {
					alreadyListed = true
				}
			}
			if !alreadyListed {
				objects = append(objects, folder)
			}
		}
	}

	return objects, nil
}
