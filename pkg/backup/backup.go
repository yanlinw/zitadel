package backup

import (
	"context"
	"github.com/caos/zitadel/pkg/backup/cockroachdb"
	"github.com/caos/zitadel/pkg/backup/rsync"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func RsyncBackupS3ToS3(
	ctx context.Context,
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationEndpoint string,
	destinationAKIDPath string,
	destinationSAKPath string,
	destinationBucket string,
	sourceName string,
	sourceEndpoint string,
	sourceAKIDPath string,
	sourceSAKPath string,
	sourceBucketPrefix string,
	configFilePath string,
) error {

	sourceAKID, err := ioutil.ReadFile(sourceAKIDPath)
	if err != nil {
		return err
	}

	sourceSAK, err := ioutil.ReadFile(sourceSAKPath)
	if err != nil {
		return err
	}

	destinationAKID, err := ioutil.ReadFile(destinationAKIDPath)
	if err != nil {
		return err
	}

	destinationSAK, err := ioutil.ReadFile(destinationSAKPath)
	if err != nil {
		return err
	}

	assetBuckets, err := ListS3AssetBuckets(sourceEndpoint, string(sourceAKID), string(sourceSAK), sourceBucketPrefix)
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartS3(sourceName, sourceEndpoint, string(sourceAKID), string(sourceSAK))
	if err != nil {
		return err
	}

	destPart, err := rsync.GetConfigPartS3(destinationName, destinationEndpoint, string(destinationAKID), string(destinationSAK))
	if err != nil {
		return err
	}

	config := strings.Join([]string{
		sourcePart,
		destPart,
	}, "\n")

	if err := ioutil.WriteFile(configFilePath, []byte(config), fs.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, assetBucket := range assetBuckets {
		wg.Add(1)
		sourceBucket := assetBucket
		if err := runCommand(
			rsync.GetCommand(
				ctx,
				configFilePath,
				sourceName,
				sourceBucket,
				destinationName,
				filepath.Join(getAssetFullPath(destinationBucket, backupName, backupNameEnv), sourceBucket),
			),
			&wg,
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func RsyncBackupS3ToGCS(
	ctx context.Context,
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationSaJsonPath string,
	destinationBucket string,
	sourceName string,
	sourceEndpoint string,
	sourceAKIDPath string,
	sourceSAKPath string,
	sourceBucketPrefix string,
	configFilePath string,
) error {
	sourceAKID, err := ioutil.ReadFile(sourceAKIDPath)
	if err != nil {
		return err
	}

	sourceSAK, err := ioutil.ReadFile(sourceSAKPath)
	if err != nil {
		return err
	}

	assetBuckets, err := ListS3AssetBuckets(sourceEndpoint, string(sourceAKID), string(sourceSAK), sourceBucketPrefix)
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartS3(sourceName, sourceEndpoint, string(sourceAKID), string(sourceSAK))
	if err != nil {
		return err
	}

	destPart, err := rsync.GetConfigPartGCS(destinationName, destinationSaJsonPath)
	if err != nil {
		return err
	}

	config := strings.Join([]string{
		sourcePart,
		destPart,
	}, "\n")

	if err := ioutil.WriteFile(configFilePath, []byte(config), fs.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, assetBucket := range assetBuckets {
		wg.Add(1)
		sourceBucket := assetBucket
		if err := runCommand(
			rsync.GetCommand(
				ctx,
				configFilePath,
				sourceName,
				sourceBucket,
				destinationName,
				filepath.Join(getAssetFullPath(destinationBucket, backupName, backupNameEnv), sourceBucket),
			),
			&wg,
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func CockroachBackupToGCS(
	ctx context.Context,
	certsFolder string,
	bucketName string,
	backupName string,
	backupNameEnv string,
	host string,
	port string,
	serviceAccountPath string,
) error {
	var wg sync.WaitGroup

	wg.Add(1)
	data, err := ioutil.ReadFile(serviceAccountPath)
	if err != nil {
		return err
	}

	err = runCommand(
		cockroachdb.GetBackupToGCS(
			ctx,
			certsFolder,
			host,
			port,
			bucketName,
			getBackupPath(backupName, backupNameEnv),
			data,
		),
		&wg,
	)
	wg.Wait()
	return err
}

func getAssetFullPath(
	bucket string,
	backupName string,
	backupNameEnv string,
) string {
	return filepath.Join(bucket, getAssetPath(backupName, backupNameEnv))
}

func getAssetPath(
	backupName string,
	backupNameEnv string,
) string {
	return filepath.Join(getBackupPath(backupName, backupNameEnv), "assets")
}

func getBackupPath(
	backupName string,
	backupNameEnv string,
) string {
	return filepath.Join(backupName, os.Getenv(backupNameEnv))
}

func CockroachBackupToS3(
	ctx context.Context,
	certsFolder string,
	bucketName string,
	backupName string,
	backupNameEnv string,
	host string,
	port string,
	accessKeyIDPath string,
	secretAccessKeyPath string,
	endpoint string,
	region string,
) error {
	var wg sync.WaitGroup

	accessKeyID, err := ioutil.ReadFile(accessKeyIDPath)
	if err != nil {
		return err
	}

	secretAccessKey, err := ioutil.ReadFile(secretAccessKeyPath)
	if err != nil {
		return err
	}

	wg.Add(1)
	err = runCommand(
		cockroachdb.GetBackupToS3(
			ctx,
			certsFolder,
			host,
			port,
			bucketName,
			getBackupPath(backupName, backupNameEnv),
			accessKeyID,
			secretAccessKey,
			[]byte(""),
			endpoint,
			region,
		),
		&wg,
	)
	wg.Wait()
	return err
}
